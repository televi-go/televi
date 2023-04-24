package connector

import (
	abstractions2 "github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/models/render"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type activeElementContext struct {
	abstractions2.TextHtmlBuilder
	keyboardBuilder *replyKbBuilder
	produceSilent   bool
	protectContent  bool
	Callbacks       *pages.Callbacks
}

func (activeElementContext *activeElementContext) BuildResult() (render.IResult, error) {

	var buttons [][]keyboards.ReplyKeyboardButton
	if activeElementContext.keyboardBuilder != nil {
		buttons = activeElementContext.keyboardBuilder.Elements
	}

	return &results.TextMessageResult{
		Text:           activeElementContext.ToString(),
		ProtectContent: activeElementContext.protectContent,
		Silent:         activeElementContext.produceSilent,
		ReplyMarkup:    results.ReplyKeyboardResult{Buttons: buttons},
	}, nil
}

func (activeElementContext *activeElementContext) ReplyKeyboard(buildAction func(builder pages.ReplyKeyboardBuilder)) {
	activeElementContext.keyboardBuilder = &replyKbBuilder{
		Callbacks:           activeElementContext.Callbacks,
		TwoDimensionBuilder: &abstractions2.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]{},
	}
	buildAction(activeElementContext.keyboardBuilder)
}
