package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/televi/core"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/external"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/core/views"
	"github.com/televi-go/televi/telegram/dto"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

type RootScene struct {
	Count    magic.State[int]
	Platform core.Platform
}

func (rootScene RootScene) Init(ctx core.InitContext) {
}

func (rootScene RootScene) Dispose() {}

func (rootScene RootScene) OnMessage(message dto.Message) {}

func (rootScene RootScene) View(builder builders.Scene) {
	builder.Head(func(headBuilder builders.Head) {
		headBuilder.SetProtection()
		media.VideoFile(headBuilder, "examples/launch_bot/sub.mov")
		headBuilder.Text("This is a televi-go bot")
		headBuilder.Row(func(builder builders.MenuRow) {
			builder.Button("Increase", func() {
				rootScene.Count.SetValueFn(func(previous int) int {
					return previous + 1
				})
			})
			builder.Button("Navigate", func() {
				rootScene.Platform.Replace(LandingScene{})
			})
		})
	})
	builder.Body(func(bodyBuilder builders.ComponentBuilder) {
		bodyBuilder.Component(views.NavigatorView(func(nav views.Navigator) builders.View {
			return BodyInnerView{Nav: nav, platform: rootScene.Platform}
		}))
	})
}

type BodyInnerViewNext struct {
	Nav      views.Navigator
	platform core.Platform
}

func (b BodyInnerViewNext) Init() {}

func (b BodyInnerViewNext) View(builder builders.ComponentBuilder) {
	builder.Message(func(viewBuilder builders.Message) {
		viewBuilder.Text("This is next")
		viewBuilder.Row(func(builder builders.ActionRowBuilder) {
			builder.Button("Go back", func(ctx builders.ClickContext) {
				b.platform.RegisterAction("", "Press back")
				b.Nav.Pop()
			})
		})
	})
}

type BodyInnerView struct {
	State     magic.State[int]
	BoldState magic.State[bool]
	Nav       views.Navigator
	platform  core.Platform
}

func (bodyInnerView BodyInnerView) Init() {}

func (bodyInnerView BodyInnerView) View(builder builders.ComponentBuilder) {
	builder.Message(func(viewBuilder builders.Message) {

		node := viewBuilder.TextF("This is sub body element\nCount is %d", bodyInnerView.State.Value())
		if bodyInnerView.BoldState.Value() {
			node.Bold()
		}
		viewBuilder.Row(func(builder builders.ActionRowBuilder) {
			builder.Button("Increase", func(ctx builders.ClickContext) {
				bodyInnerView.platform.RegisterAction("", "Press increase")
				bodyInnerView.State.SetValueFn(func(previous int) int {
					return previous + 1
				})
			})
			builder.Button("Transit", func(ctx builders.ClickContext) {
				bodyInnerView.Nav.Push(BodyInnerViewNext{Nav: bodyInnerView.Nav, platform: bodyInnerView.platform})
			})
			if bodyInnerView.BoldState.Value() {
				builder.Button("Make regular", func(ctx builders.ClickContext) {
					bodyInnerView.BoldState.SetValue(false)
				})
			} else {
				builder.Button("MAKE BOLD", func(ctx builders.ClickContext) {
					bodyInnerView.BoldState.SetValue(true)
				})
			}
		})
	})
}

type LandingScene struct {
	TabViewState magic.State[int]
	TimerState   magic.State[time.Duration]
}

func (landing LandingScene) View(builder builders.Scene) {
	builder.Head(func(head builders.Head) {
		head.Text("Ваши клиенты вас заждались!")
		head.Row(func(builder builders.MenuRow) {
			builder.Button("Opt1", func() {
				landing.TabViewState.SetValue(0)
			})
			builder.Button("Opt2", func() {
				landing.TabViewState.SetValue(1)
			})
		})
	})
	builder.Body(func(body builders.ComponentBuilder) {
		body.Message(func(message builders.Message) {
			message.TextF("Timer value: %v\n", landing.TimerState.Value())
		})
		body.Message(func(viewBuilder builders.Message) {
			if landing.TabViewState.Value() == 0 {
				viewBuilder.Text("First variant")
			}
			if landing.TabViewState.Value() == 1 {
				viewBuilder.Text("Second variant")
			}
		})
	})
}

func (landing LandingScene) Init(ctx core.InitContext) {
	ctx.OnExternal("timer", func(data any) {
		landing.TimerState.SetValue(data.(time.Duration))
	})
}

func (landing LandingScene) Dispose() {

}

func (landing LandingScene) OnMessage(message dto.Message) {}

func main() {
	fmt.Printf("Process %d\nEnv: %s\n", os.Getpid(), strings.Join(os.Environ(), "\n"))
	app, err := core.NewAppBuilder().
		WithEnvToken("Token").
		WithDefaultAddress().
		WithRootScene(func(platform core.Platform) core.ActionScene {
			return RootScene{Platform: platform}
		}).
		WithServerSetup(func(router fiber.Router, db *sql.DB) {

		}).
		Build()
	if err != nil {
		log.Fatalln(err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		beginTime := time.Now()
		for t := range ticker.C {
			app.DispatchExternal("timer", external.AllUsersTarget{}, t.Sub(beginTime))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()
	app.Run(ctx)
}
