package connector

import (
	"televi/connector/abstractions"
	"televi/models/pages"
	"televi/models/render"
	"televi/models/render/results"
	"televi/telegram/messages/keyboards"
)

type textElementContext struct {
	abstractions.TextHtmlBuilder
	componentPrefix string
	keyboardBuilder *inlineKeyboardBuilder
	produceSilent   bool
	protectContent  bool
	callbacks       *pages.Callbacks
}

func (textElement *textElementContext) ButtonsRow(builder func(rowBuilder pages.InlineRowBuilder)) {
	if textElement.keyboardBuilder == nil {
		textElement.keyboardBuilder = &inlineKeyboardBuilder{
			Callbacks:         textElement.callbacks,
			CbComponentPrefix: textElement.componentPrefix,
		}
	}

	textElement.keyboardBuilder.ButtonsRow(builder)
}

func (textElement *textElementContext) BuildResult() (render.IResult, error) {

	var buttons [][]keyboards.InlineKeyboardButton
	if textElement.keyboardBuilder != nil && textElement.keyboardBuilder.Elements != nil {
		buttons = textElement.keyboardBuilder.Elements
	}

	return &results.TextMessageResult{
		Text:           textElement.ToString(),
		ProtectContent: textElement.protectContent,
		Silent:         textElement.produceSilent,
		ReplyMarkup:    &results.InlineKeyboardResult{Keyboard: buttons},
	}, nil
}

func (textElement *textElementContext) InlineKeyboard(buildAction func(builder pages.InlineKeyboardBuilder)) {
	kbBuilder := &inlineKeyboardBuilder{
		Callbacks:         textElement.callbacks,
		CbComponentPrefix: textElement.componentPrefix,
	}
	buildAction(kbBuilder)
	textElement.keyboardBuilder = kbBuilder
}
