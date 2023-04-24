package results

import (
	"github.com/televi-go/televi/models/render"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"io"
)

type SingleMediaResult struct {
	Text           string
	ProtectContent bool
	Silent         bool
	Key            string
	FileId         string
	FileName       string
	FileReader     io.Reader
	ReplyMarkup    KeyboardResult
	Type           string
	HasSpoiler     bool
}

func (photoResult *SingleMediaResult) Kind() string {
	return "single" + photoResult.Type
}

func (photoResult *SingleMediaResult) InitAction(destination telegram.Destination) telegram.Request {
	return messages.SingleMediaRequest{
		MediaType: photoResult.Type,
		Base: messages.MediaMessageBase{
			Destination:    destination,
			Caption:        photoResult.Text,
			ProtectContent: photoResult.ProtectContent,
			Silent:         photoResult.Silent,
			ReplyTo:        0,
			ReplyMarkup:    photoResult.ReplyMarkup.ToReplyMarkup(),
		},
		FileName:    photoResult.FileName,
		Content:     nil, //photoResult.FileReader,
		PhotoFileId: photoResult.FileId,
		HasSpoiler:  photoResult.HasSpoiler,
	}
}

func (photoResult *SingleMediaResult) CompareTo(result render.IResult, destination telegram.Destination, messageIds []int) (canBeUpdated bool, requests []telegram.Request) {
	photoResultNext, isPhotoResult := result.(*SingleMediaResult)
	if !isPhotoResult {
		return false, nil
	}

	kbAction := photoResult.ReplyMarkup.CanBeUpdated(photoResultNext.ReplyMarkup)

	switch kbAction {
	case NoAction:
		break
	case EditAction:
		// Assume that ReplyKeyboard never outputs EditAction from comparison
		requests = append(requests, keyboards.EditInlineKeyboardRequest{
			Destination: destination,
			MessageId:   messageIds[0],
			NewKeyboard: photoResultNext.ReplyMarkup.ToReplyMarkup(),
		})
	case ReplaceAction:
		return false, nil
	}

	if photoResultNext.Text != photoResult.Text {
		requests = append(requests, messages.EditMessageCaptionRequest{
			EditInlineKeyboardRequest: keyboards.EditInlineKeyboardRequest{
				Destination: destination,
				MessageId:   messageIds[0],
				NewKeyboard: photoResultNext.ReplyMarkup.ToReplyMarkup(),
			},
			Caption: photoResultNext.Text,
		})
	}

	if photoResultNext.Key != photoResult.Key {
		requests = append(requests, messages.UpdateMediaRequest{
			EditMessageCaptionRequest: messages.EditMessageCaptionRequest{
				EditInlineKeyboardRequest: keyboards.EditInlineKeyboardRequest{
					Destination: destination,
					MessageId:   messageIds[0],
					NewKeyboard: photoResultNext.ReplyMarkup.ToReplyMarkup(),
				},
				Caption: photoResultNext.Text,
			},
			Media: messages.InputMedia{
				Raw:     photoResultNext.FileReader,
				FileId:  photoResultNext.FileId,
				Type:    photoResultNext.Type,
				Caption: photoResultNext.Text,
			},
		})
	}

	if len(requests) > 0 {
		requests = requests[len(requests)-1:]
	}

	return true, requests
}
