package connector

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type replyKbBuilder struct {
	*abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]
	Callbacks *pages.Callbacks
}

func (replyKeyboardBuilder *replyKbBuilder) ActionButton(caption string, callback pages.ClickCallback) {
	replyKeyboardBuilder.Callbacks.AddMessageListener(func(message *dto.Message) {

		if message.Text == caption {
			callback()
		}
	})
	replyKeyboardBuilder.Add(keyboards.ReplyKeyboardButton{Text: caption})
}

func (replyKeyboardBuilder *replyKbBuilder) ContactButton(caption string, callback pages.ContactCallback) {
	replyKeyboardBuilder.Callbacks.AddMessageListener(func(msg *dto.Message) {
		callback(msg.Contact)
	})
	replyKeyboardBuilder.Add(keyboards.ReplyKeyboardButton{
		Text:            caption,
		RequestContact:  true,
		RequestLocation: false,
	})
}

func (replyKeyboardBuilder *replyKbBuilder) ButtonsRow(builder func(rowBuilder pages.ReplyRowBuilder)) {
	builder(replyKeyboardBuilder)
	replyKeyboardBuilder.CommitRow()
}
