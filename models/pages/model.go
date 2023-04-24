package pages

import (
	"github.com/televi-go/televi/models"
	"github.com/televi-go/televi/telegram/bot"
)

type TransitionKind int

const (
	SeparativeTransition TransitionKind = iota
	ReplacingTransition  TransitionKind = iota
)

type Model struct {
	Page            Scene
	Result          *models.ResultLine
	Previous        *Model
	Callbacks       Callbacks
	Kind            TransitionKind
	Origin          ViewSequenceOrigin
	BoundRespondIds []int
}

type ViewSequenceNode struct {
	View      Scene
	Previous  *ViewSequenceNode
	Callbacks Callbacks
	Result    *models.ResultLine
}

type ViewSequence struct {
	Current *ViewSequenceNode
	Origin  ViewSequenceOrigin
}

type ViewSequenceOrigin interface {
	Remove(api *bot.Api)
}
