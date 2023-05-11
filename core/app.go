package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/televi-go/televi/core/external"
	"github.com/televi-go/televi/core/metrics"
	"github.com/televi-go/televi/profiler"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"log"
	"os"
	"strings"
	"sync"
)

type App struct {
	controllers      map[string]*Controller
	controllerAccess sync.Mutex
	api              *bot.Api
	initScene        func(Platform) ActionScene
	context          context.Context
	Profiler         *profiler.Throughput
	Server           *metrics.ServerImpl
}

func (app *App) getUserList() map[string]bool {
	m := make(map[string]bool, len(app.controllers))
	for dest := range app.controllers {
		m[dest] = true
	}
	return m
}

func (app *App) DispatchExternal(event string, to external.Target, data any) {
	targets := to.GetFromUserList(app.getUserList())
	for _, target := range targets {
		app.controllers[target].Dispatch(ExternalEvent{Domain: &DomainEvent{
			Name: event,
			Data: data,
		}})
	}
}

func (app *App) getUpdates() {
	updateChannel := app.api.Poll(app.context)
	for update := range updateChannel {
		app.dispatchUpdate(update)
	}
}

func NewApp(
	token string,
	address string,
	initScene func(platform Platform) ActionScene,
	serverImpl *metrics.ServerImpl,
) (*App, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	if address == "" {
		return nil, errors.New("address is empty")
	}

	api := bot.NewApi(token, address)
	return &App{
		controllers:      map[string]*Controller{},
		controllerAccess: sync.Mutex{},
		api:              api,
		initScene:        initScene,
		context:          nil,
		Server:           serverImpl,
		Profiler: &profiler.Throughput{
			LabelPrefix: "bot@" + strings.Split(token, ":")[0],
		},
	}, nil
}

func (app *App) getOrCreateController(destination telegram.Destination, info *dto.User) *Controller {
	app.controllerAccess.Lock()
	defer app.controllerAccess.Unlock()
	controller, hasController := app.controllers[destination.ToString()]
	if !hasController {
		controller = NewController(destination, app.api, app.initScene, info, app.Profiler, app.Server.DB)

		if app.Server != nil {
			err := metrics.AddRegistered(app.Server.DB, info)
			if err != nil {
				log.Printf("error inserting new user: %v\n", err)
			}
		}

		app.controllers[destination.ToString()] = controller
		go controller.Run(app.context)
	}
	return controller
}
func (app *App) dispatchUpdate(update dto.Update) {
	controller := app.getOrCreateController(telegram.ChatDestination{ChatId: int(update.SentFrom().ID)}, update.SentFrom())
	controller.Dispatch(ExternalEvent{
		Message:  update.Message,
		Callback: update.CallbackQuery,
	})
}

func (app *App) acquirePidLock() (func(), error) {
	appName := fmt.Sprintf("/var/run/televi/%d.pid", app.api.GetBotId())
	currentProcess, noFileErr := os.ReadFile(appName)
	if noFileErr == nil {
		return nil, fmt.Errorf("another process %s is running this bot (%d)", currentProcess, app.api.GetBotId())
	}
	writeErr := os.WriteFile(appName, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if writeErr != nil {
		log.Printf("Unable to acquire pid lock for bot %d\n. One instance is not guaranteed", app.api.GetBotId())
		return func() {}, nil
	}
	return func() {

		os.Remove(appName)
	}, nil
}

func (app *App) Run(ctx context.Context) {
	lockDestructor, err := app.acquirePidLock()
	if err != nil {
		log.Printf("abort: %v", err)
		return
	}
	defer lockDestructor()
	defer app.Profiler.WriteStats(os.Stdout)

	if app.Server != nil {
		log.Printf("launching server\n")
		go func() {
			err := app.Server.Serve(ctx, "connect.sock")
			if err != nil {
				log.Printf("err %v\n", err)
			}
		}()
	}

	app.context = ctx
	app.getUpdates()
}
