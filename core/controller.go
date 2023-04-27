package core

import (
	"context"
	"github.com/televi-go/televi/core/body"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/profiler"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages"
)

type Controller struct {
	navStack                  NavigationStack
	destination               telegram.Destination
	stateUpdateChan           chan body.NeedsPaintTask
	navigationChan            chan NavigationTask
	externalCommunicationChan chan ExternalEvent
	api                       *bot.Api
	context                   context.Context
	userInfo                  *dto.User
	Profiler                  *profiler.Throughput
}

func NewController(destination telegram.Destination, api *bot.Api, initSceneCtor func(platform Platform) ActionScene, info *dto.User, profiler *profiler.Throughput) *Controller {
	navC := make(chan NavigationTask, 10)
	initScene := initSceneCtor(platformImpl{
		NavImpl: NavImpl{
			NavC: navC,
		},
		user: *info,
	})
	stateUpdateChannel := make(chan body.NeedsPaintTask, 10)
	magic.InjectInPlace(initScene, func() {
		stateUpdateChannel <- body.NeedsPaintTask{SceneStateUpdated: true}
	})

	controller := &Controller{
		navStack: NavigationStack{
			Current:     NewNavStackEntry(initScene, stateUpdateChannel, destination, nil),
			Destination: destination,
		},
		destination:               destination,
		stateUpdateChan:           stateUpdateChannel,
		navigationChan:            navC,
		externalCommunicationChan: make(chan ExternalEvent, 100),
		api:                       api,
		userInfo:                  info,
		Profiler:                  profiler,
	}
	controller.paint(true)
	return controller
}

func (controller *Controller) processEvent(event ExternalEvent) {
	if event.Message != nil {
		controller.processMessage(event.Message)
	}
	if event.Callback != nil {
		controller.processCallback(event.Callback)
	}

	if event.Domain != nil {
		controller.currentStackEntry().ExternalCallbacks.Dispatch(event.Domain.Name, event.Domain.Data)
	}
}

func (controller *Controller) processCallback(callback *dto.CallbackQuery) {
	ctx := body.ClickContextImpl{AnswerRequest: messages.AnswerCallbackRequest{
		Id:        callback.ID,
		Text:      "",
		ShowAlert: false,
	}}
	controller.currentStackEntry().BodyCallbacks.Execute(callback.Data, &ctx)
	controller.api.LaunchRequest(ctx.AnswerRequest)
}

func (controller *Controller) processMessage(message *dto.Message) {
	callback, hasCallback := controller.currentStackEntry().HeadResult.HeadCallbacks.OnRegularButton[message.Text]
	controller.currentStackEntry().BodyResult.AddInterruption(message.MessageID)
	if hasCallback {
		callback()
		return
	}
	controller.currentStackEntry().ActionScene.OnMessage(*message)
}

func (controller *Controller) currentStackEntry() *NavStackEntry {
	return controller.navStack.Current
}

func (controller *Controller) currentScene() ActionScene {
	return controller.currentStackEntry().ActionScene
}

func (controller *Controller) paint(needsBodyRepaint bool) {
	sw := controller.Profiler.NewStopwatch("paint")
	defer sw.Record()
	controller.currentStackEntry().Render(needsBodyRepaint, controller.api)
}

func (controller *Controller) Dispatch(event ExternalEvent) {
	controller.externalCommunicationChan <- event
}

func (controller *Controller) processNavTask(task NavigationTask) {
	defer controller.paint(true)
	switch task.Action {
	case ExtendAction:
		controller.navStack.push(task.Target, controller.stateUpdateChan, true, controller.api)
		break
	case ReplaceAction:
		controller.navStack.push(task.Target, controller.stateUpdateChan, false, controller.api)
		break
	case PopAction:
		controller.navStack.pop()
		break
	case PopToMainAction:
		controller.navStack.popMain()
		break
	}
}

type platformImpl struct {
	NavImpl
	user dto.User
}

func (p platformImpl) GetUser() dto.User {
	return p.user
}

func (controller *Controller) Run(ctx context.Context) {
	defer close(controller.externalCommunicationChan)
	defer close(controller.navigationChan)
	defer close(controller.stateUpdateChan)
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-controller.externalCommunicationChan:
			go controller.processEvent(evt)
			break
		case task := <-controller.navigationChan:
			controller.processNavTask(task)
			break
		case stateUpdate := <-controller.stateUpdateChan:
			controller.paint(stateUpdate.SceneStateUpdated)
			break
		}
	}
}

type ExternalEvent struct {
	Domain   *DomainEvent
	Message  *dto.Message
	Callback *dto.CallbackQuery
}

type DomainEvent struct {
	Name string
	Data any
}

type MessageEvent struct {
	Message *dto.Message
}
