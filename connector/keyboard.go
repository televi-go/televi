package connector

import (
	"fmt"
	"gtihub.com/televi-go/televi/connector/abstractions"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/telegram/messages/keyboards"
)

type inlineKeyboardBuilder struct {
	abstractions.TwoDimensionBuilder[keyboards.InlineKeyboardButton]
	Callbacks         *pages.Callbacks
	CbComponentPrefix string
}

func (inlineKeyBoardBuilder *inlineKeyboardBuilder) ActionButton(caption string, callback pages.ClickCallback) {
	cbData := fmt.Sprintf("%s:%s", inlineKeyBoardBuilder.CbComponentPrefix, caption)
	inlineKeyBoardBuilder.Callbacks.AddButtonListener(pages.EventData{
		Kind:    "",
		Payload: cbData,
	}, callback)
	inlineKeyBoardBuilder.Add(keyboards.InlineKeyboardButton{
		Caption:      caption,
		Url:          "",
		CallbackData: cbData,
		WebApp:       keyboards.WebAppInfo{},
	})
}

func (inlineKeyBoardBuilder *inlineKeyboardBuilder) UrlButton(caption string, url string) {
	inlineKeyBoardBuilder.Add(keyboards.InlineKeyboardButton{
		Caption:      caption,
		Url:          url,
		CallbackData: "",
		WebApp:       keyboards.WebAppInfo{},
	})
}

func (inlineKeyBoardBuilder *inlineKeyboardBuilder) ButtonsRow(builder func(rowBuilder pages.InlineRowBuilder)) {
	builder(inlineKeyBoardBuilder)
	inlineKeyBoardBuilder.CommitRow()
}
