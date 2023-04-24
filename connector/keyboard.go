package connector

import (
	"fmt"
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type inlineKeyboardBuilder struct {
	abstractions.TwoDimensionBuilder[keyboards.InlineKeyboardButton]
	Callbacks         *pages.Callbacks
	CbComponentPrefix string
}

func (inlineKeyBoardBuilder *inlineKeyboardBuilder) BuildKeyboard() (result results.KeyboardResult, err error) {
	var buttons [][]keyboards.InlineKeyboardButton
	if inlineKeyBoardBuilder != nil && inlineKeyBoardBuilder.Elements != nil {
		buttons = inlineKeyBoardBuilder.Elements
	}
	return &results.InlineKeyboardResult{Keyboard: buttons}, nil
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
