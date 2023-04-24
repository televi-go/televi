package messages

import (
	"errors"
	"github.com/televi-go/televi/telegram"
)

type DeleteMessageRequest struct {
	MessageId   int
	Destination telegram.Destination
}

func (deleteMessageRequest DeleteMessageRequest) Method() string {
	return "deleteMessage"
}

func (deleteMessageRequest DeleteMessageRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)

	if deleteMessageRequest.MessageId == 0 {
		return nil, errors.New("unspecified message to delete")
	}

	params.WriteNonZero("message_id", deleteMessageRequest.MessageId)
	err := deleteMessageRequest.Destination.WriteParameter(params)
	return params, err
}
