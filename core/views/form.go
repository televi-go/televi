package views

import (
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/telegram/dto"
)

type formView[TForm any] struct {
	empty          TForm
	current        magic.State[TForm]
	partProcessors []*processorBind[TForm]
	onFill         func(TForm)
}

type processorBind[TForm any] struct {
	proc         FormPartProcessor[TForm]
	intro        builders.View
	hasSucceeded bool
}

type FormPartProcessor[TForm any] func(msg dto.Message, form *TForm) error

func (f formView[TForm]) Init() {

}

func (f formView[TForm]) View(builder builders.ComponentBuilder) {

}
