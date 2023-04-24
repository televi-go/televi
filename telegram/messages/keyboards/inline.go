package keyboards

import (
	"github.com/televi-go/televi/telegram"
)

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (inlineKeyboard InlineKeyboardMarkup) replyMarkupImpl() {}

func (inlineKeyboard InlineKeyboardMarkup) WriteParameter(params telegram.Params) error {
	if inlineKeyboard.Keyboard == nil {
		inlineKeyboard.Keyboard = make([][]InlineKeyboardButton, 0)
	}
	return params.WriteJson("reply_markup", inlineKeyboard)
}

type WebAppInfo struct {
	Url string `json:"url"`
}

type InlineKeyboardButton struct {
	Caption      string     `json:"text"`
	Url          string     `json:"url,omitempty"`
	CallbackData string     `json:"callback_data,omitempty"`
	WebApp       WebAppInfo `json:"web_app,omitempty"`
}
