package views

import (
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
)

type navigator struct {
	stack   []builders.View
	Trigger magic.Trigger
}

func (navigator *navigator) Push(fragment builders.View) {
	navigator.stack = append(navigator.stack, fragment)
	navigator.Trigger.SendChange()
}

func (navigator *navigator) Pop() {
	if len(navigator.stack) > 1 {
		navigator.stack = navigator.stack[:len(navigator.stack)-1]
		navigator.Trigger.SendChange()
	}
}

func (navigator *navigator) Init() {}

func (navigator *navigator) View(builder builders.ComponentBuilder) {
	builder.Component(navigator.stack[len(navigator.stack)-1])
}

type Navigator interface {
	Push(fragment builders.View)
	Pop()
}

func NavigatorView(rootView func(nav Navigator) builders.View) builders.View {
	nav := &navigator{}
	rootViewImpl := rootView(nav)
	nav.stack = append(nav.stack, rootViewImpl)
	return nav
}
