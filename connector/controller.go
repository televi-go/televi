package connector

import (
	"context"
	"fmt"
	"gtihub.com/televi-go/televi/delayed"
	"gtihub.com/televi-go/televi/models"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/models/render"
	"gtihub.com/televi-go/televi/telegram"
	"gtihub.com/televi-go/televi/telegram/bot"
	"gtihub.com/televi-go/televi/telegram/dto"
	"gtihub.com/televi-go/televi/telegram/messages"
	"log"

	"time"
)

type Controller struct {
	ChatId                 telegram.Destination
	CurrentModel           *pages.Model
	Api                    *bot.Api
	EventChannel           chan ControllerReactionEvent
	Scheduler              *delayed.TaskScheduler
	LatestRender           time.Time
	StateUpdateChannel     chan struct{}
	ActiveCallbacks        pages.Callbacks
	isCurrentlyRehydrating bool
	ActiveCallbacksModel   *pages.Model
}

type TransitionOptions struct {
	TransitionPage   pages.Scene
	TransitPolicy    pages.TransitPolicy
	TransitBack      bool
	TransitionOrigin pages.ViewSequenceOrigin
	TransitToMain    bool
	OnPerformed      chan bool
}

type ControllerReactionEvent struct {
	TelegramMessage  *dto.Message
	TelegramCallback *dto.CallbackQuery
	//TODO: handle other update types
	InnerStateChanged bool
	ExternalEvent     *models.ExternalEvent
	IsRehydrate       bool
	transitOptions    *TransitionOptions
}

type MessageTransitionOrigin struct {
	destination telegram.Destination
	messageId   int
}

func (origin MessageTransitionOrigin) Remove(api *bot.Api) {
	go api.Request(messages.DeleteMessageRequest{
		MessageId:   origin.messageId,
		Destination: origin.destination,
	})
}

const rehydrateDelta = time.Hour * 47

func (controller *Controller) dispatchMessage(message *dto.Message) bool {
	reactCtx := &reactContextImpl{
		controller:             controller,
		message:                message,
		includeTransitToAnchor: controller.ActiveCallbacksModel != controller.CurrentModel,
		origin: MessageTransitionOrigin{
			destination: controller.ChatId,
			messageId:   message.MessageID,
		},
	}
	hasCallbacks := controller.ActiveCallbacks.InvokeButton(pages.EventData{
		Kind: "message-reply", Payload: message.Text}, reactCtx)
	if hasCallbacks {
		return true
	}
	return controller.CurrentModel.Callbacks.InvokeOnMessage(reactCtx)
}

func (controller *Controller) dispatchCallback(callback *dto.CallbackQuery) bool {
	reactCtx := &reactContextImpl{
		controller: controller,
		message:    nil,
		AlertRequest: messages.AnswerCallbackRequest{
			Id:        callback.ID,
			Text:      "",
			ShowAlert: false,
		},
	}
	hasCalled := controller.CurrentModel.Callbacks.InvokeButton(pages.EventData{
		Kind:    "",
		Payload: callback.Data,
	}, reactCtx)

	go func() {
		resp, err := controller.Api.Request(reactCtx.AlertRequest)
		if err != nil {
			fmt.Println(resp, err)
		}
	}()

	return hasCalled && !reactCtx.WasTransitRequested
}

func (controller *Controller) dispatchExternal(event *models.ExternalEvent) bool {
	panic("External events not implemented yet")
}

func (controller *Controller) dispatchEvent(event ControllerReactionEvent) {
	hasChanged := controller.processEvent(event)

	if event.IsRehydrate {
		controller.render(true)
		return
	}

	if !hasChanged && !event.InnerStateChanged {
		return
	}
}

func (controller *Controller) _transitBack() bool {
	if controller.CurrentModel.Previous == nil {
		return false
	}
	if controller.CurrentModel.Origin != nil {
		controller.CurrentModel.Origin.Remove(controller.Api)
	}
	if controller.CurrentModel.Kind == pages.SeparativeTransition {
		for _, completedResult := range controller.CurrentModel.Result.Line {
			for _, cleanupRequest := range completedResult.Cleanup(controller.ChatId) {
				go controller.Api.Request(cleanupRequest)
			}
		}
	}

	controller.CurrentModel = controller.CurrentModel.Previous
	return true
}

func (controller *Controller) dispatchTransition(options TransitionOptions) (result bool) {

	defer func() {
		fmt.Printf("transited to %T\n", controller.CurrentModel.Page)
		controller.StateUpdateChannel <- struct{}{}
	}()
	if options.TransitionPage != nil {

		resultLine := &render.ResultLine{}
		pages.MountStates(&options.TransitionPage, controller.StateUpdateChannel)
		if options.TransitPolicy.KeepPrevious {
			resultLine = controller.CurrentModel.Result
		}

		controller.CurrentModel = &pages.Model{
			Page:      options.TransitionPage,
			Result:    resultLine,
			Origin:    options.TransitionOrigin,
			Previous:  controller.CurrentModel,
			Callbacks: *pages.NewCallbacks(),
			Kind:      options.TransitPolicy.GetKind(),
		}
		return true
	}

	if options.TransitBack {
		return controller._transitBack()
	}

	if options.TransitToMain {
		for controller._transitBack() {
		}
		return true
	}

	return false
}

func (controller *Controller) processEvent(event ControllerReactionEvent) (hasCalled bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error in handling event: ", r)
			hasCalled = true
		}
	}()
	if event.TelegramMessage != nil {
		hasCalled = controller.dispatchMessage(event.TelegramMessage)
	}
	if event.TelegramCallback != nil {
		hasCalled = controller.dispatchCallback(event.TelegramCallback)
	}
	if event.ExternalEvent != nil {
		hasCalled = controller.dispatchExternal(event.ExternalEvent)
	}

	if event.transitOptions != nil {
		controller.dispatchTransition(*event.transitOptions)
		return true
	}

	return
}

func (controller *Controller) render(silent bool) {
	fmt.Printf("rendering %T\n", controller.CurrentModel.Page)
	if silent {
		fmt.Println("Rehydration called")
	}

	chatId, _ := controller.ChatId.(telegram.ChatDestination)

	pageBuildContext := &BuildContext{
		elements:        nil,
		everySilent:     silent,
		everyProtected:  false,
		ActiveCallbacks: pages.NewCallbacks(),
		Callbacks:       pages.NewCallbacks(),
		UserId:          chatId.ChatId,
	}
	controller.CurrentModel.Page.View(pageBuildContext)
	controller.CurrentModel.Callbacks = *pageBuildContext.Callbacks
	if !pageBuildContext.ActiveCallbacks.IsEmpty() {
		controller.ActiveCallbacks = *pageBuildContext.ActiveCallbacks
		controller.ActiveCallbacksModel = controller.CurrentModel
	}

	newLine := pageBuildContext.buildLine()
	result := controller.CurrentModel.Result.CompareAndProduce(controller.ChatId, newLine, silent)
	err := controller.CurrentModel.Result.Run(result, controller.Api)
	if err != nil {
		fmt.Println("error in running: ", err)
	}
	controller.LatestRender = time.Now()
	err = controller.Scheduler.Schedule(
		"rehydrate",
		time.Now().Add(rehydrateDelta),
		controller.ChatId.ToString(),
	)
	//fmt.Println("planned rehydrate at", time.Now().Add(rehydrateDelta))
	if err != nil {
		log.Println("Error in rehydration planning", err)
	}
}

func (controller *Controller) EnqueueRehydrate() {

	if controller.LatestRender.Add(rehydrateDelta).After(time.Now()) {
		return
	}

	controller.Enqueue(ControllerReactionEvent{
		TelegramMessage:   nil,
		TelegramCallback:  nil,
		InnerStateChanged: false,
		ExternalEvent:     nil,
		IsRehydrate:       true,
	})
}

func (controller *Controller) Enqueue(event ControllerReactionEvent) {
	controller.EventChannel <- event
}

func (controller *Controller) RunQueue(ctx context.Context) {
	defer close(controller.StateUpdateChannel)
	defer close(controller.EventChannel)
	for {
		select {
		case <-ctx.Done():
			return
		case <-controller.StateUpdateChannel:
			controller.render(false)
			break
		case evt := <-controller.EventChannel:
			controller.dispatchEvent(evt)
			break
		}
	}
}

func (controller *Controller) transitTo(page pages.Scene, policy pages.TransitPolicy, origin pages.ViewSequenceOrigin) {
	controller.Enqueue(ControllerReactionEvent{transitOptions: &TransitionOptions{
		TransitionPage:   page,
		TransitPolicy:    policy,
		TransitBack:      false,
		TransitionOrigin: origin,
	}})

}

func (controller *Controller) transitBack() {
	controller.Enqueue(ControllerReactionEvent{
		transitOptions: &TransitionOptions{TransitBack: true},
	})

}

func (controller *Controller) transitToMain() {
	controller.Enqueue(ControllerReactionEvent{
		transitOptions: &TransitionOptions{TransitToMain: true},
	})

}
