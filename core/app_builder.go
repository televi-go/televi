package core

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/televi/core/metrics"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"os"
	"strconv"
)

type appBuilder struct {
	token      string
	address    string
	rootScene  func(platform Platform) ActionScene
	serverImpl *metrics.ServerImpl
}

func (builder appBuilder) WithServerSetup(f func(router fiber.Router, db *sql.DB)) Builder {
	port, _ := strconv.Atoi(os.Getenv("DEBUG_PORT"))
	db, err := sql.Open("mysql", "root:@/televi?parseTime=true")
	if err != nil {
		panic(err)
	}

	api := bot.NewApi(builder.token, builder.address)
	resp, err := api.Request(bot.GetMeRequest{})
	if err != nil {
		panic(fmt.Errorf("error establishing bot %v, %s", err, builder.token))
	}

	user, _ := telegram.ParseAs[dto.User](resp)

	builder.serverImpl = &metrics.ServerImpl{
		Port:     port,
		DB:       db,
		FiberApp: fiber.New(),
		Info:     user.UserName,
	}

	err = builder.serverImpl.Setup(f)
	if err != nil {
		panic(err)
	}

	return builder
}

func NewAppBuilder() AppTokenConsumer {
	return appBuilder{}
}

func (builder appBuilder) WithToken(token string) ApiAddressConsumer {
	return appBuilder{
		token:     token,
		address:   "",
		rootScene: builder.rootScene,
	}
}

func (builder appBuilder) WithEnvToken(env string) ApiAddressConsumer {
	return appBuilder{
		token:     os.Getenv(env),
		address:   "",
		rootScene: builder.rootScene,
	}
}

func (builder appBuilder) WithApiAddress(address string) RootSceneConsumer {
	return appBuilder{
		token:     builder.token,
		address:   address,
		rootScene: builder.rootScene,
	}
}

func (builder appBuilder) WithDefaultAddress() RootSceneConsumer {
	return builder.WithApiAddress("https://api.telegram.org")
}

func (builder appBuilder) WithRootScene(f func(platform Platform) ActionScene) Builder {
	return appBuilder{address: builder.address, token: builder.token, rootScene: f}
}

func (builder appBuilder) Build() (*App, error) {
	return NewApp(builder.token, builder.address, builder.rootScene, builder.serverImpl)
}

type Builder interface {
	Build() (*App, error)
	WithServerSetup(func(router fiber.Router, db *sql.DB)) Builder
}

type AppTokenConsumer interface {
	WithToken(token string) ApiAddressConsumer
	WithEnvToken(env string) ApiAddressConsumer
}

type RootSceneConsumer interface {
	WithRootScene(func(platform Platform) ActionScene) Builder
}

type ApiAddressConsumer interface {
	WithApiAddress(address string) RootSceneConsumer
	WithDefaultAddress() RootSceneConsumer
}
