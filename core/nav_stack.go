package core

import (
	"github.com/televi-go/televi/core/body"
	"github.com/televi-go/televi/core/head_management"
	"github.com/televi-go/televi/telegram"
)

type NavStackEntry struct {
	ActionScene    ActionScene
	WasInitialized bool
	Previous       *NavStackEntry
	Kind           NavigationKind
	HeadResult     *head_management.HeadResultContainer
	BodyRoot       *body.Root
	BodyResult     *body.Result
	BodyCallbacks  *body.Callbacks
}

type NavigationStack struct {
	Current     *NavStackEntry
	Destination telegram.Destination
}

type NavigationKind int

const (
	ReplaceNavigationKind NavigationKind = iota
	ExtendNavigationKind  NavigationKind = iota
)
