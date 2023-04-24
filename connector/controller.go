package connector

import (
	"context"
	"fmt"
	"github.com/televi-go/televi/delayed"
	"github.com/televi-go/televi/models"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages"
	"github.com/televi-go/televi/util"
	"log"

	"time"
)

type NavigationProvider interface {
	EnqueueTransit(options TransitionOptions)
	dispatchAlert(show bool, text string)
}

type Controller struct {
	ChatId                 telegram.Destination
	CurrentModel           *pages.Model
	Api                    *bot.Api
	UserInfo               *dto.User
	EventChannel           chan ControllerReactionEvent
	Scheduler              *delayed.TaskScheduler
	LatestRender           time.Time
	StateUpdateChannel     chan struct{}
	isCurrentlyRehydrating bool
	ReplyCallbacks         *pages.Callbacks
	LastMessage            *dto.Message
	LastQuery              *dto.CallbackQuery
}

type TransitionOptions struct {
	From        *pages.Model
	To          pages.Scene
	IsExtending bool
	IsBack      bool
	IsToMain    bool
	IsReplace   bool
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
	defer func() {
		if controller.LastMessage != nil {
			controller.CurrentModel.BoundRespondIds = append(controller.CurrentModel.BoundRespondIds, controller.LastMessage.MessageID)
		}
	}()
	controller.UserInfo = message.From
	controller.LastMessage = message

	canConsumeWithKeyboard := controller.ReplyCallbacks.InvokeOnMessage(message)
	canConsumeElse := controller.CurrentModel.Callbacks.InvokeOnMessage(message)
	return canConsumeWithKeyboard || canConsumeElse
}

func (controller *Controller) dispatchAlert(showAlert bool, text string) {
	if controller.LastQuery != nil {
		query := controller.LastQuery
		controller.LastQuery = nil
		request := messages.AnswerCallbackRequest{
			Id:        query.ID,
			Text:      text,
			ShowAlert: showAlert,
		}
		go func() {
			resp, err := controller.Api.Request(request)
			if err != nil {
				fmt.Println(resp, err)
			}
		}()
	}
}

func (controller *Controller) dispatchCallback(callback *dto.CallbackQuery) bool {
	controller.UserInfo = callback.From
	hasCalled := controller.CurrentModel.Callbacks.InvokeButton(pages.EventData{
		Kind:    "",
		Payload: callback.Data,
	})

	/* go func() {
		resp, err := controller.Api.Request(AlertRequest)
		if err != nil {
			fmt.Println(resp, err)
		}
	}() */
	controller.dispatchAlert(false, "")
	return hasCalled
}

func (controller *Controller) dispatchExternal(event *models.ExternalEvent) bool {
	panic("External events not implemented yet")
}

func (controller *Controller) dispatchEvent(event ControllerReactionEvent) {
	controller.processEvent(event)

	if event.IsRehydrate {
		controller.render(true)
		return
	}
}

func (controller *Controller) _transitBack(remove bool) bool {
	if controller.CurrentModel.Previous == nil {
		return false
	}
	if controller.CurrentModel.Origin != nil && remove {
		controller.CurrentModel.Origin.Remove(controller.Api)
	}
	if controller.CurrentModel.Kind == pages.SeparativeTransition && remove {
		for _, completedResult := range controller.CurrentModel.Result.Line {
			for _, cleanupRequest := range completedResult.Cleanup(controller.ChatId) {
				go controller.Api.Request(cleanupRequest)
			}
		}
		for _, id := range controller.CurrentModel.BoundRespondIds {
			go controller.Api.Request(messages.DeleteMessageRequest{MessageId: id, Destination: controller.ChatId})
		}
	}

	controller.CurrentModel = controller.CurrentModel.Previous
	return true
}

func (controller *Controller) dispatchTransition(options TransitionOptions) (result bool) {
	defer func() {
		if result {
			controller.StateUpdateChannel <- struct{}{}
		}
	}()
	if options.IsToMain {
		hasTransited := false
		for controller.CurrentModel.Previous != nil {
			hasTransited = controller._transitBack(!options.IsReplace)
		}
		return hasTransited
	}

	if options.IsBack {
		for controller.CurrentModel != options.From && controller.CurrentModel.Previous != nil {
			controller._transitBack(!options.IsReplace)
		}

		return controller._transitBack(!options.IsReplace)
	}

	var origin pages.ViewSequenceOrigin

	if controller.LastMessage != nil {
		origin = MessageTransitionOrigin{
			destination: controller.ChatId,
			messageId:   controller.LastMessage.MessageID,
		}
		controller.LastMessage = nil
	}

	controller.transit(options.From, options.To, options.IsExtending, origin, options.IsReplace)

	return true
}

func (controller *Controller) EnqueueTransit(options TransitionOptions) {
	controller.Enqueue(ControllerReactionEvent{transitOptions: &options})
}

func (controller *Controller) transit(
	from *pages.Model,
	to pages.Scene,
	isExtending bool,
	origin pages.ViewSequenceOrigin,
	isReplacing bool,
) {

	for controller.CurrentModel != from && controller.CurrentModel.Previous != nil {
		controller._transitBack(!isReplacing)
	}

	resultLine := &models.ResultLine{}
	pages.MountStates(&to, controller.StateUpdateChannel)
	if !isExtending && !isReplacing {
		resultLine = controller.CurrentModel.Result
	}

	var prev = controller.CurrentModel

	if isReplacing {
		prev = controller.CurrentModel.Previous
	}

	controller.CurrentModel = &pages.Model{
		Page:      to,
		Result:    resultLine,
		Origin:    origin,
		Previous:  prev,
		Callbacks: *pages.NewCallbacks(),
		Kind:      pages.TransitPolicy{KeepPrevious: !isExtending}.GetKind(),
	}
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
	sw := util.NewStopWatch()
	defer func() {
		fmt.Println("Painted in", sw.Elapsed())
	}()

	if silent {
		fmt.Println("Rehydration called")
	}

	//chatId, _ := controller.ChatId.(telegram.ChatDestination)

	pageBuildContext := &BuildContext{
		UserInfo:        controller.UserInfo,
		elements:        nil,
		everySilent:     silent,
		everyProtected:  false,
		ActiveCallbacks: pages.NewCallbacks(),
		Callbacks:       pages.NewCallbacks(),
		controller:      controller,
		stackPoint:      controller.CurrentModel,
	}

	controller.CurrentModel.Page.View(pageBuildContext)
	controller.CurrentModel.Callbacks = *pageBuildContext.Callbacks

	if !pageBuildContext.ActiveCallbacks.IsEmpty() {
		controller.ReplyCallbacks = pageBuildContext.ActiveCallbacks
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
