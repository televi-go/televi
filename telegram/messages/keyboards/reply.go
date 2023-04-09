package keyboards

import "televi/telegram"

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
}

func (replyKeyboardRemove ReplyKeyboardRemove) WriteParameter(params telegram.Params) error {
	params.WriteBool("remove_keyboard", true)
	return nil
}

func (replyKeyboardRemove ReplyKeyboardRemove) replyMarkupImpl() {}
