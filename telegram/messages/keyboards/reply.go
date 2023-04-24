package keyboards

import "github.com/televi-go/televi/telegram"

type ReplyKeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]ReplyKeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                    `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool                    `json:"one_time_keyboard,omitempty"`
	Selective       bool                    `json:"selective,omitempty"`
}

func (replyKeyboard ReplyKeyboardMarkup) replyMarkupImpl() {}

func (replyKeyboard ReplyKeyboardMarkup) WriteParameter(params telegram.Params) error {
	return params.WriteJson("reply_markup", replyKeyboard)
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
}

func (replyKeyboardRemove ReplyKeyboardRemove) WriteParameter(params telegram.Params) error {
	replyKeyboardRemove.RemoveKeyboard = true
	return params.WriteJson("reply_markup", replyKeyboardRemove)
}

func (replyKeyboardRemove ReplyKeyboardRemove) replyMarkupImpl() {}
