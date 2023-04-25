package main

import (
	"context"
	"fmt"
	"github.com/televi-go/televi/core"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/core/runner"
	"github.com/televi-go/televi/core/views"
	"github.com/televi-go/televi/telegram/dto"
	"log"
	"os"
	"os/signal"
)

type RootScene struct {
	Count magic.State[int]
}

func (rootScene RootScene) Init() {

}

func (rootScene RootScene) Dispose() {}

func (rootScene RootScene) OnMessage(message dto.Message) {}

func (rootScene RootScene) View(builder builders.Scene) {
	builder.Head(func(headBuilder builders.Head) {
		media.ImageFile(headBuilder, "examples/launch_bot/welcome_pic.png")
		headBuilder.Text("This is a televi-go bot")
		headBuilder.Row(func(builder builders.MenuRow) {
			builder.Button("Increase", func() {
				rootScene.Count.SetValueFn(func(previous int) int {
					return previous + 1
				})
			})
		})
	})
	builder.Body(func(bodyBuilder builders.ComponentBuilder) {
		bodyBuilder.Component(views.NavigatorView(func(nav views.Navigator) builders.View {
			return BodyInnerView{Nav: nav}
		}))
	})
}

type BodyInnerViewNext struct {
	Nav views.Navigator
}

func (b BodyInnerViewNext) Init() {}

func (b BodyInnerViewNext) View(builder builders.ComponentBuilder) {
	builder.Message(func(viewBuilder builders.Message) {
		viewBuilder.Text("This is next")
		viewBuilder.Row(func(builder builders.ActionRowBuilder) {
			builder.Button("Go back", func(ctx builders.ClickContext) {
				b.Nav.Pop()
			})
		})
	})
}

type BodyInnerView struct {
	State     magic.State[int]
	BoldState magic.State[bool]
	Nav       views.Navigator
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
				bodyInnerView.State.SetValueFn(func(previous int) int {
					return previous + 1
				})
			})
			builder.Button("Transit", func(ctx builders.ClickContext) {
				bodyInnerView.Nav.Push(BodyInnerViewNext{Nav: bodyInnerView.Nav})
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

func main() {
	fmt.Printf("Process %d\n", os.Getpid())
	app, err := runner.NewApp(
		os.Getenv("Token"),
		"https://api.telegram.org",
		func() core.ActionScene {
			return RootScene{}
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()
	app.Run(ctx)
}
