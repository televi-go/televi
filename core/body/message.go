package body

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type RealMessageBuilder interface {
	builders.Message
	SetIsModified(value bool)
}

func (message Message) FirstMedia() media.Media {
	return message.Media[0]
}

type Message struct {
	Media          []media.Media
	Text           string
	Actions        results.InlineKeyboardResult
	IsCached       bool
	ProtectContent bool

	// TODO: poll!!
}

type MessageKind int

const (
	TextKind        MessageKind = iota
	SingleMediaKind MessageKind = iota
	MediaGroupKind  MessageKind = iota
	PollKind        MessageKind = iota
)

func (message Message) GetKind() MessageKind {
	switch len(message.Media) {
	case 0:
		return TextKind
	case 1:
		return SingleMediaKind
	default:
		return MediaGroupKind
	}
}

func (message Message) InitRequest(destination telegram.Destination) telegram.Request {
	switch message.GetKind() {
	case TextKind:
		return messages.TextMessageRequest{
			Destination:    destination,
			Text:           message.Text,
			ProtectContent: message.ProtectContent,
			Silent:         false,
			ReplyTo:        0,
			ReplyMarkup:    message.Actions.ToReplyMarkup(),
		}
	case SingleMediaKind:
		return messages.SingleMediaRequest{
			Base: messages.MediaMessageBase{
				Destination:    destination,
				Caption:        message.Text,
				ProtectContent: message.ProtectContent,
				Silent:         false,
				ReplyTo:        0,
				ReplyMarkup:    message.Actions.ToReplyMarkup(),
			},
			Content:     message.Media[0].Content,
			FileName:    "",
			PhotoFileId: message.Media[0].FileId,
			HasSpoiler:  message.Media[0].HasSpoiler,
			MediaType:   message.Media[0].FieldName(),
		}
	}
	panic("invalid kind")
}

type ViewMounter interface {
	Add(view builders.InMessageView)
}

type MessageProducer interface {
	Build() Message
}

type MessageBuilderImpl struct {
	media       []media.Media
	Callbacks   *Callbacks
	IsProtected bool
	abstractions.TextHtmlBuilder
	abstractions.TwoDimensionBuilder[keyboards.InlineKeyboardButton]
	IsCached bool
}

func (m *MessageBuilderImpl) SetIsModified(value bool) {
	m.IsCached = !value
}

func (m *MessageBuilderImpl) SetProtection() {
	m.IsProtected = true
}

func (m *MessageBuilderImpl) Build() Message {
	return Message{
		Media: m.media,
		Text:  m.ToString(),
		Actions: results.InlineKeyboardResult{
			Keyboard: m.Elements,
		},
		IsCached:       m.IsCached,
		ProtectContent: m.IsProtected,
	}
}

func (m *MessageBuilderImpl) Button(caption string, onclick func(ctx builders.ClickContext)) {
	data := m.Callbacks.bind(caption, onclick)
	m.Add(keyboards.InlineKeyboardButton{
		Caption:      caption,
		Url:          "",
		CallbackData: data,
	})
}

func (m *MessageBuilderImpl) Url(caption string, target string) {
	m.Add(keyboards.InlineKeyboardButton{
		Caption:      caption,
		Url:          target,
		CallbackData: "",
		WebApp:       keyboards.WebAppInfo{},
	})
}

func (m *MessageBuilderImpl) Row(f func(builder builders.ActionRowBuilder)) {
	f(m)
	m.CommitRow()
}

func (m *MessageBuilderImpl) Media(media media.Media) {
	m.media = append(m.media, media)
}

func (m *MessageBuilderImpl) Component(view builders.InMessageView) {
	panic("unreachable code")
}
