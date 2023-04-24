package head_management

import (
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
)

type HeadResult struct {
	Media          *media.Media
	Text           string
	ProtectContent bool
	Keyboard       results.ReplyKeyboardResult
}

func (result HeadResult) equals(other HeadResult) bool {
	return result.Text == other.Text &&
		result.MediaKind() == other.MediaKind() &&
		KeyboardsAreSame(result.Keyboard, other.Keyboard) &&
		result.mediaKey() == other.mediaKey()
}

func (result HeadResult) mediaKey() string {
	if result.Media == nil {
		return ""
	}
	return result.Media.Key
}

func (result HeadResult) MediaKind() media.Kind {
	if result.Media == nil {
		return media.NoMedia
	}
	return result.Media.Kind
}

func (result HeadResult) InitRequest(destination telegram.Destination) telegram.Request {
	if result.Media == nil {
		return messages.TextMessageRequest{
			Destination:    destination,
			Text:           result.Text,
			ProtectContent: result.ProtectContent,
			Silent:         false,
			ReplyTo:        0,
			ReplyMarkup:    result.Keyboard.ToReplyMarkup(),
		}
	}
	return messages.SingleMediaRequest{
		Base: messages.MediaMessageBase{
			Destination:    destination,
			Caption:        result.Text,
			ProtectContent: result.ProtectContent,
			Silent:         false,
			ReplyTo:        0,
			ReplyMarkup:    result.Keyboard.ToReplyMarkup(),
		},
		Content:     result.Media.Content,
		FileName:    "",
		PhotoFileId: result.Media.FileId,
		HasSpoiler:  result.Media.HasSpoiler,
		MediaType:   result.Media.FieldName(),
	}
}
