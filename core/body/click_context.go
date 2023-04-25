package body

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/telegram/messages"
)

type ClickContextImpl struct {
	AnswerRequest messages.AnswerCallbackRequest
}

func (clickCtx *ClickContextImpl) Alert(builder func(builders.AlertBuilder)) {
	alertBuilder := AlertBuilder{}
	builder(&alertBuilder)
	clickCtx.AnswerRequest.Text = alertBuilder.ToString()
	clickCtx.AnswerRequest.ShowAlert = alertBuilder.Mode == builders.AlertMode
}

type AlertBuilder struct {
	abstractions.TextHtmlBuilder
	Mode builders.DisplayMode
}

func (a *AlertBuilder) SetDisplayMode(mode builders.DisplayMode) {
	a.Mode = mode
}
