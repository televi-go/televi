package connector

import (
	"gtihub.com/televi-go/televi/connector/abstractions"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/telegram/dto"
	"gtihub.com/televi-go/televi/telegram/messages/keyboards"
)

type replyKbBuilder struct {
	*abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]
	Callbacks *pages.Callbacks
}

func (replyKeyboardBuilder *replyKbBuilder) ActionButton(caption string, callback pages.ClickCallback) {
	replyKeyboardBuilder.Callbacks.AddButtonListener(pages.EventData{
		Kind:    "message-reply",
		Payload: caption,
	}, callback)
	replyKeyboardBuilder.Add(keyboards.ReplyKeyboardButton{Text: caption})
}

type contactContext struct {
	pages.MessageReactionContext
}

func (c contactContext) Contact() *dto.Contact {
	return c.Message().Contact
}

func (replyKeyboardBuilder *replyKbBuilder) ContactButton(caption string, callback pages.ContactCallback) {
	replyKeyboardBuilder.Callbacks.AddMessageListener(func(ctx pages.MessageReactionContext) {
		contactCtx := contactContext{ctx}
		callback(contactCtx)
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
