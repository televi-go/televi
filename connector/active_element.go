package connector

import (
	"televi/connector/abstractions"
	"televi/models/pages"
	"televi/models/render"
	"televi/models/render/results"
	"televi/telegram/messages/keyboards"
)

type activeElementContext struct {
	abstractions.TextHtmlBuilder
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
		TwoDimensionBuilder: &abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]{},
	}
	buildAction(activeElementContext.keyboardBuilder)
}
