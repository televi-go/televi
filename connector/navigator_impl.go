package connector

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/models/pages"
)

type navigator struct {
	controller NavigationProvider
	stackPoint *pages.Model
}

func (navigator navigator) Alert(f func(alertBuilder pages.TextPartBuilder)) {
	textBuilder := &abstractions.TextHtmlBuilder{}
	f(textBuilder)
	navigator.controller.dispatchAlert(true, textBuilder.ToString())
}

func (navigator navigator) Push(page pages.Scene) {
	navigator.controller.EnqueueTransit(TransitionOptions{
		From:        navigator.stackPoint,
		To:          page,
		IsExtending: false,
		IsBack:      false,
		IsToMain:    false,
	})
}

func (navigator navigator) Extend(page pages.Scene) {
	navigator.controller.EnqueueTransit(TransitionOptions{
		From:        navigator.stackPoint,
		To:          page,
		IsExtending: true,
		IsBack:      false,
		IsToMain:    false,
	})
}

func (navigator navigator) Replace(page pages.Scene) {
	navigator.controller.EnqueueTransit(TransitionOptions{
		From:        navigator.stackPoint,
		To:          page,
		IsExtending: false,
		IsBack:      false,
		IsToMain:    false,
		IsReplace:   true,
	})
}

func (navigator navigator) Pop() {
	navigator.controller.EnqueueTransit(TransitionOptions{
		From:        navigator.stackPoint,
		To:          nil,
		IsExtending: false,
		IsBack:      true,
		IsToMain:    false,
	})
}

func (navigator navigator) PopAll() {
	navigator.controller.EnqueueTransit(TransitionOptions{IsToMain: true})
}
