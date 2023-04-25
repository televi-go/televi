package core

import (
	"context"
	"github.com/televi-go/televi/core/body"
	"github.com/televi-go/televi/core/head_management"
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
	navigationChan            chan struct{}
	externalCommunicationChan chan ExternalEvent
	api                       *bot.Api
	context                   context.Context
	userInfo                  *dto.User
	Profiler                  *profiler.Throughput
}

func NewController(destination telegram.Destination, api *bot.Api, initScene ActionScene, info *dto.User, profiler *profiler.Throughput) *Controller {
	stateUpdateChannel := make(chan body.NeedsPaintTask, 10)
	magic.InjectInPlace(initScene, func() {
		stateUpdateChannel <- body.NeedsPaintTask{SceneStateUpdated: true}
	})
	controller := &Controller{
		navStack: NavigationStack{
			Current: &NavStackEntry{
				BodyCallbacks: body.NewCallbacks(),
				ActionScene:   initScene,
				Previous:      nil,
				Kind:          0,
				HeadResult:    &head_management.HeadResultContainer{},
				BodyRoot:      body.NewRoot(stateUpdateChannel),
				BodyResult:    &body.Result{},
			},
			Destination: nil,
		},
		destination:               destination,
		stateUpdateChan:           stateUpdateChannel,
		navigationChan:            make(chan struct{}, 10),
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
	if hasCallback {
		callback()
	}
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
	bodyRoot := controller.currentStackEntry().BodyRoot
	if needsBodyRepaint {
		bodyRoot.SetNeedsUpdate()
	}
	sceneContext := NewSceneContext(bodyRoot)
	controller.currentScene().View(sceneContext)
	var wasReplaced = false
	if needsBodyRepaint {
		headResult := sceneContext.HeadBuilder.Build()
		headResultContainer := controller.currentStackEntry().HeadResult
		wasReplaced = headResultContainer.CompareAgainst(headResult, controller.destination, controller.api)
		headResultContainer.HeadCallbacks = sceneContext.HeadBuilder.Callbacks
	}

	bodyResultLine, callbacks := sceneContext.Root.ProvideResult()
	controller.currentStackEntry().BodyResult.CompareAgainst(
		bodyResultLine,
		controller.api,
		controller.destination,
		wasReplaced,
	)
	controller.currentStackEntry().BodyCallbacks = callbacks
}

func (controller *Controller) Dispatch(event ExternalEvent) {
	controller.externalCommunicationChan <- event
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
			controller.processEvent(evt)
		case <-controller.navigationChan:
			controller.paint(true)
		case stateUpdate := <-controller.stateUpdateChan:
			controller.paint(stateUpdate.SceneStateUpdated)
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
