package main

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/televi/core"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/telegram/dto"
	"log"
	"os"
	"os/signal"
)

type CredForm struct {
	_start   string
	Name     string
	LastName string
	Address  string
}

func (form CredForm) CompleteField(field string) CredForm {

	if form._start == "" {
		return CredForm{_start: field}
	}

	if form.Name == "" {
		return CredForm{_start: form._start, Name: field}
	}
	if form.LastName == "" {
		return CredForm{_start: form._start, Name: form.Name, LastName: field}
	}
	return CredForm{_start: form._start, Name: form.Name, LastName: form.LastName, Address: field}
}

type FormScene struct {
	platform core.Platform
	State    magic.State[CredForm]
}

func (f FormScene) View(builder builders.Scene) {
	builder.Head(func(head builders.Head) {
		head.Text("Please fill the form")
	})
	builder.Body(func(body builders.ComponentBuilder) {
		body.Message(func(message builders.Message) {
			message.Text("Please fill your firstname")
		})
		if f.State.Value().Name != "" {
			body.Message(func(message builders.Message) {
				message.Text("Please fill your last name")
			})
		}
		if f.State.Value().LastName != "" {
			body.Message(func(message builders.Message) {
				message.Text("Please fill your address")
			})
		}
	})
}

func (f FormScene) Init(ctx core.InitContext) {

}

func (f FormScene) Dispose() {
}

func (f FormScene) OnMessage(message dto.Message) {
	f.State.SetValueFn(func(previous CredForm) CredForm {
		return previous.CompleteField(message.Text)
	})
}

func main() {
	app, err := core.NewAppBuilder().
		WithEnvToken("Token").
		WithDefaultAddress().
		WithRootScene(func(platform core.Platform) core.ActionScene {
			return FormScene{platform: platform}
		}).WithServerSetup(func(router fiber.Router, db *sql.DB) {

	}).Build()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	app.Run(ctx)
}
