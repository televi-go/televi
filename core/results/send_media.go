package results

import (
	"github.com/televi-go/televi/core/update"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type SendMediaRequest struct {
	MediaType string
	Content   []byte
	Text      string
	Keyboard  keyboards.ReplyMarkup
}

func (s SendMediaRequest) GetRequest(destination telegram.Destination) telegram.Request {
	return messages.SingleMediaRequest{
		Base: messages.MediaMessageBase{
			Destination:    destination,
			Caption:        s.Text,
			ProtectContent: false,
			Silent:         false,
			ReplyTo:        0,
			ReplyMarkup:    s.Keyboard,
		},
		Content:     s.Content,
		FileName:    "",
		PhotoFileId: "",
		HasSpoiler:  false,
		MediaType:   s.MediaType,
	}
}

func (s SendMediaRequest) InflateResult(response telegram.Response) update.Update {
	return nil
}
