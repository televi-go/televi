package keyboards

import "github.com/televi-go/televi/telegram"

type ReplyMarkup interface {
	telegram.ParamsWriter
	replyMarkupImpl()
}

type EditInlineKeyboardRequest struct {
	Destination telegram.Destination
	MessageId   int
	NewKeyboard ReplyMarkup
}

func (editKeyboardRequest EditInlineKeyboardRequest) Method() string {
	return "editMessageReplyMarkup"
}

func (editKeyboardRequest EditInlineKeyboardRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)
	var err error
	_ = editKeyboardRequest.Destination.WriteParameter(params)
	params.WriteInt("message_id", editKeyboardRequest.MessageId)
	if editKeyboardRequest.NewKeyboard != nil {
		err = editKeyboardRequest.NewKeyboard.WriteParameter(params)
	}
	return params, err
}
