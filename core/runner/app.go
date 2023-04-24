package runner

import (
	"context"
	"errors"
	"github.com/televi-go/televi/core"
	"github.com/televi-go/televi/profiler"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"os"
	"strings"
	"sync"
)

type App struct {
	controllers      map[string]*core.Controller
	controllerAccess sync.Mutex
	api              *bot.Api
	initScene        func() core.ActionScene
	context          context.Context
	Profiler         *profiler.Throughput
}

func (app *App) getUpdates() {
	updateChannel := app.api.Poll(app.context)
	for update := range updateChannel {
		app.dispatchUpdate(update)
	}
}

func NewApp(token string, address string, initScene func() core.ActionScene) (*App, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	if address == "" {
		return nil, errors.New("address is empty")
	}

	api := bot.NewApi(token, address)

	return &App{
		controllers:      map[string]*core.Controller{},
		controllerAccess: sync.Mutex{},
		api:              api,
		initScene:        initScene,
		context:          nil,
		Profiler: &profiler.Throughput{
			LabelPrefix: "bot@" + strings.Split(token, ":")[0],
		},
	}, nil
}

func (app *App) getOrCreateController(destination telegram.Destination, info *dto.User) *core.Controller {
	app.controllerAccess.Lock()
	defer app.controllerAccess.Unlock()
	controller, hasController := app.controllers[destination.ToString()]
	if !hasController {
		controller = core.NewController(destination, app.api, app.initScene(), info, app.Profiler)
		app.controllers[destination.ToString()] = controller
		go controller.Run(app.context)
	}
	return controller
}
func (app *App) dispatchUpdate(update dto.Update) {
	controller := app.getOrCreateController(telegram.ChatDestination{ChatId: int(update.SentFrom().ID)}, update.SentFrom())
	controller.Dispatch(core.ExternalEvent{
		Message:  update.Message,
		Callback: update.CallbackQuery,
	})
}

func (app *App) Run(ctx context.Context) {
	defer app.Profiler.WriteStats(os.Stdout)
	app.context = ctx
	app.getUpdates()
}
