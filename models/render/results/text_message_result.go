package results

import (
	"televi/models/render"
	"televi/telegram"
	"televi/telegram/messages"
	"televi/telegram/messages/keyboards"
	"televi/util"
)

type TextMessageResult struct {
	Text           string
	ProtectContent bool
	Silent         bool
	ReplyMarkup    KeyboardResult
}

func (textMessageResult *TextMessageResult) Kind() string {
	return "text"
}

func (textMessageResult *TextMessageResult) InitAction(destination telegram.Destination) telegram.Request {
	return messages.TextMessageRequest{
		Destination:    destination,
		Text:           textMessageResult.Text,
		ProtectContent: textMessageResult.ProtectContent,
		Silent:         textMessageResult.Silent,
		ReplyMarkup:    textMessageResult.ReplyMarkup.ToReplyMarkup(),
	}
}

func (textMessageResult *TextMessageResult) CompareTo(result render.IResult, destination telegram.Destination, messageIds []int) (bool, []telegram.Request) {
	tmResult, isTmResult := result.(*TextMessageResult)
	if !isTmResult {
		return false, nil
	}

	keyboardAction := textMessageResult.ReplyMarkup.CanBeUpdated(tmResult.ReplyMarkup)

	if keyboardAction == ReplaceAction {
		return false, nil
	}

	if keyboardAction == EditAction || tmResult.Text != textMessageResult.Text {

		rm := tmResult.ReplyMarkup.ToReplyMarkup()
		_, isNotInlineErr := util.PointerOr[keyboards.InlineKeyboardMarkup](rm)
		if isNotInlineErr != nil {
			return false, nil
		}
		return true, []telegram.Request{
			messages.EditMessageRequest{
				Destination: destination,
				Text:        tmResult.Text,
				MessageId:   messageIds[0],
				ReplyMarkup: rm,
			},
		}
	}

	return true, nil
}
