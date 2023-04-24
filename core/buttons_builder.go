package core

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/callbacks"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type ButtonsBuilder struct {
	abstractions.TwoDimensionBuilder[InlineButton]
	Provider callbacks.ViewCallbacksProvider
}

func (b *ButtonsBuilder) Button(caption string, onclick func()) {
	b.Provider.Bind(caption, onclick)
	b.Add(InlineButton{
		Caption: caption,
		Url:     "",
	})
}

func (b *ButtonsBuilder) Url(caption string, target string) {
	b.Add(InlineButton{
		Caption: caption,
		Url:     target,
	})
}

func (b *ButtonsBuilder) Row(f func(builder builders.ActionRowBuilder)) {
	b.CommitRow()
	f(b)
	b.CommitRow()
}

type InlineButton struct {
	Caption string
	Url     string
}

type KeyboardBuilder struct {
	abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]
	//Callbacks *callbacks.SceneCallbacks
}
