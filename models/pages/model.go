package pages

import (
	"gtihub.com/televi-go/televi/models/render"
	"gtihub.com/televi-go/televi/telegram/bot"
)

type TransitionKind int

const (
	SeparativeTransition TransitionKind = iota
	ReplacingTransition  TransitionKind = iota
)

type Model struct {
	Page      Scene
	Result    *render.ResultLine
	Previous  *Model
	Callbacks Callbacks
	Kind      TransitionKind
	Origin    ViewSequenceOrigin
}

type ViewSequenceNode struct {
	View      Scene
	Previous  *ViewSequenceNode
	Callbacks Callbacks
	Result    *render.ResultLine
}

type ViewSequence struct {
	Current *ViewSequenceNode
	Origin  ViewSequenceOrigin
}

type ViewSequenceOrigin interface {
	Remove(api *bot.Api)
}
