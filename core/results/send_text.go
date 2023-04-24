package results

import (
	"github.com/televi-go/televi/core/update"
	results2 "github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
)

type SendText struct {
	Text     string
	Keyboard results2.ReplyKeyboardResult
}

func (s SendText) GetRequest(destination telegram.Destination) telegram.Request {
	return messages.TextMessageRequest{
		Destination:    destination,
		Text:           s.Text,
		ProtectContent: false,
		Silent:         false,
		ReplyTo:        0,
		ReplyMarkup:    s.Keyboard.ToReplyMarkup(),
	}
}

func (s SendText) InflateResult(response telegram.Response) update.Update {
	return s
}
