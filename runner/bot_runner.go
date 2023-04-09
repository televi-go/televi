package runner

import (
	"context"
	"fmt"
	"gtihub.com/televi-go/televi/connector"
	"gtihub.com/televi-go/televi/delayed"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/models/render"
	"gtihub.com/televi-go/televi/telegram"
	"gtihub.com/televi-go/televi/telegram/bot"
	"gtihub.com/televi-go/televi/telegram/dto"
)

type Runner struct {
	controllers     map[string]*connector.Controller
	primaryPageCtor func() pages.Scene
	api             *bot.Api
	ctx             context.Context
	scheduler       *delayed.TaskScheduler
}

func (runner *Runner) getUpdates() {
	updateChannel := runner.api.Poll(runner.ctx)
	for update := range updateChannel {
		runner.dispatchUpdate(update)
	}
}

func (runner *Runner) getOrCreateController(destination telegram.Destination) *connector.Controller {
	controller, hasController := runner.controllers[destination.ToString()]
	if !hasController {
		fmt.Println("Creating controller for", destination)
		stateUpdateChannel := make(chan struct{}, 1)
		page := runner.primaryPageCtor()
		pages.MountStates(&page, stateUpdateChannel)
		controller = &connector.Controller{
			ChatId:          destination,
			ActiveCallbacks: *pages.NewCallbacks(),
			CurrentModel: &pages.Model{
				Page: page,
				Result: &render.ResultLine{
					Line: nil,
				},
				Previous:  nil,
				Callbacks: *pages.NewCallbacks(),
			},
			Api:                runner.api,
			StateUpdateChannel: stateUpdateChannel,
			EventChannel:       make(chan connector.ControllerReactionEvent, 10),
			Scheduler:          runner.scheduler,
		}
		runner.controllers[destination.ToString()] = controller
		// nothing -> some state

		go func() {
			controller.RunQueue(runner.ctx)
		}()
		controller.StateUpdateChannel <- struct{}{}
	}
	return controller
}

func (runner *Runner) dispatchUpdate(update dto.Update) {
	destination := telegram.GetDestination(update)
	controller := runner.getOrCreateController(destination)
	controller.Enqueue(connector.ControllerReactionEvent{
		TelegramMessage:   update.Message,
		TelegramCallback:  update.CallbackQuery,
		InnerStateChanged: false,
		ExternalEvent:     nil,
	})
}

func (runner *Runner) Run(ctx context.Context) {
	runner.ctx = ctx
	go func() {
		runner.scheduler.Run(ctx)
	}()
	runner.getUpdates()
}

const DefaultAPiAddress = "https://api.telegram.org"

func NewRunner(token string, ctor func() pages.Scene, dsn string, address string) (*Runner, error) {
	scheduler, err := delayed.NewScheduler(dsn)
	if err != nil {
		return nil, err
	}
	runner := &Runner{
		controllers:     map[string]*connector.Controller{},
		primaryPageCtor: ctor,
		api:             bot.NewApi(token, address),
		ctx:             nil,
		scheduler:       scheduler,
	}

	delayed.Register(scheduler, "rehydrate", func(args string) {
		controller := runner.getOrCreateController(telegram.ParseDestination(args))
		controller.EnqueueRehydrate()
	})

	return runner, err
}
