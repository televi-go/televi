package main

import (
	"github.com/televi-go/televi/core"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/telegram/dto"
)

type RootScene struct {
	FileId magic.State[string]
	Kind   magic.State[string]
}

func (r RootScene) View(builder builders.Scene) {
	builder.Head(func(head builders.Head) {
		head.Text("Send me media and i will output fileids")
	})

	builder.Body(func(body builders.ComponentBuilder) {
		if r.FileId.Value() != "" {
			body.Message(func(message builders.Message) {
				message.TextF("This is %s. FileId = %s", r.Kind.Value(), r.FileId.Value())
			})
		}
	})

}

func (r RootScene) Init(ctx core.InitContext) {

}

func (r RootScene) Dispose() {

}

func (r RootScene) OnMessage(message dto.Message) {
	if message.Video != nil {
		r.Kind.SetValue("video")
		r.FileId.SetValue(message.Video.FileID)
	}

	if message.Sticker != nil {
		r.Kind.SetValue("sticker")
		r.FileId.SetValue(message.Sticker.FileID)
	}

	if message.VideoNote != nil {
		r.Kind.SetValue("video note")
		r.FileId.SetValue(message.VideoNote.FileID)
	}
}
