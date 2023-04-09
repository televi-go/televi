package messages

import (
	"gtihub.com/televi-go/televi/telegram"
	"gtihub.com/televi-go/televi/telegram/messages/keyboards"
)

type TextMessageRequest struct {
	Destination    telegram.Destination
	Text           string
	ProtectContent bool
	Silent         bool
	ReplyTo        int
	ReplyMarkup    keyboards.ReplyMarkup
}

func (textMessageRequest TextMessageRequest) Method() string {
	return "sendMessage"
}

func (textMessageRequest TextMessageRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)
	err := textMessageRequest.Destination.WriteParameter(params)
	if err != nil {
		return nil, err
	}

	params.WriteString("text", textMessageRequest.Text)
	params.WriteString("parse_mode", "HTML")
	params.WriteBool("protect_content", textMessageRequest.ProtectContent)
	params.WriteBool("disable_notification", textMessageRequest.Silent)
	params.WriteNonZero("reply_to_message_id", textMessageRequest.ReplyTo)
	if textMessageRequest.ReplyMarkup != nil {
		err = textMessageRequest.ReplyMarkup.WriteParameter(params)
	}

	return params, err
}

type EditMessageRequest struct {
	Destination telegram.Destination
	Text        string
	MessageId   int
	ReplyMarkup keyboards.ReplyMarkup
}

func (editMessageRequest EditMessageRequest) Method() string {
	return "editMessageText"
}

func (editMessageRequest EditMessageRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)

	err := editMessageRequest.Destination.WriteParameter(params)
	if err != nil {
		return nil, err
	}

	params.WriteNonZero("message_id", editMessageRequest.MessageId)
	params.WriteString("text", editMessageRequest.Text)
	params.WriteString("parse_mode", "HTML")
	if editMessageRequest.ReplyMarkup != nil {
		err = editMessageRequest.ReplyMarkup.WriteParameter(params)
	}
	return params, err
}
