package pages

import (
	"televi/models/render"
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
}
