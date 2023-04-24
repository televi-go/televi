package body

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type FragmentBuilderImpl struct {
	//UpdateC  chan<- struct{}
	Callbacks *Callbacks
	Messages  []Message
}

func (fragmentBuilder *FragmentBuilderImpl) Component(view builders.View) {
	panic("unreachable code")
}

func (fragmentBuilder *FragmentBuilderImpl) Message(builder func(viewBuilder builders.Message)) {
	fragmentBuilder.Callbacks.beginNextMessage()
	viewBuilder := &MessageBuilderImpl{
		media:               nil,
		Callbacks:           fragmentBuilder.Callbacks,
		TextHtmlBuilder:     abstractions.TextHtmlBuilder{},
		TwoDimensionBuilder: abstractions.TwoDimensionBuilder[keyboards.InlineKeyboardButton]{},
	}
	builder(viewBuilder)
	fragmentBuilder.Messages = append(fragmentBuilder.Messages, viewBuilder.Build())
}
