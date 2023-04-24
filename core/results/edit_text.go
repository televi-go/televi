package results

import (
	"github.com/televi-go/televi/core/update"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
)

type EditMessageText struct {
	Text      string
	MessageId int
}

func (e EditMessageText) GetRequest(destination telegram.Destination) telegram.Request {
	return messages.EditMessageRequest{
		Destination: destination,
		Text:        "",
		MessageId:   e.MessageId,
		ReplyMarkup: nil,
	}
}

func (e EditMessageText) InflateResult(response telegram.Response) update.Update {
	return e
}
