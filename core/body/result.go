package body

import (
	"bytes"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"github.com/televi-go/televi/util"
	"log"
	"runtime"
)

type ResultEntry struct {
	Message
	BoundIds []int
}

type Result struct {
	Entries []*ResultEntry
}

func (entry *ResultEntry) cleanup(
	api *bot.Api,
	destination telegram.Destination,
) {
	for _, messageId := range entry.BoundIds {
		if messageId != 0 {
			api.LaunchRequest(messages.DeleteMessageRequest{
				MessageId:   messageId,
				Destination: destination,
			})
		}
	}
}

func (entry *ResultEntry) compareAsText(
	newer Message, destination telegram.Destination) telegram.Request {
	return messages.EditMessageRequest{
		Destination: destination,
		Text:        newer.Text,
		MessageId:   entry.BoundIds[0],
		ReplyMarkup: newer.Actions.ToReplyMarkup(),
	}
}

func (entry *ResultEntry) compareAsSingleMedia(
	newer Message,
	destination telegram.Destination,
) telegram.Request {
	if entry.FirstMedia().Key != newer.FirstMedia().Key {
		return messages.UpdateMediaRequest{
			EditMessageCaptionRequest: messages.EditMessageCaptionRequest{
				EditInlineKeyboardRequest: keyboards.EditInlineKeyboardRequest{
					Destination: destination,
					MessageId:   entry.BoundIds[0],
					NewKeyboard: newer.Actions.ToReplyMarkup(),
				},
				Caption: newer.Text,
			},
			Media: messages.InputMedia{
				Raw:     bytes.NewReader(newer.FirstMedia().Content),
				FileId:  newer.FirstMedia().FileId,
				Type:    newer.FirstMedia().FieldName(),
				Caption: newer.Text,
			},
		}
	}
	if entry.Text != newer.Text {
		return messages.EditMessageCaptionRequest{
			EditInlineKeyboardRequest: keyboards.EditInlineKeyboardRequest{
				Destination: destination,
				MessageId:   entry.BoundIds[0],
				NewKeyboard: newer.Actions.ToReplyMarkup(),
			},
			Caption: newer.Text,
		}
	}
	return keyboards.EditInlineKeyboardRequest{
		Destination: destination,
		MessageId:   entry.BoundIds[0],
		NewKeyboard: newer.Actions.ToReplyMarkup(),
	}
}

func (entry *ResultEntry) compareNonReplace(
	newer Message,
	destination telegram.Destination,
) (result []telegram.Request) {

	switch entry.Message.GetKind() {
	case TextKind:
		return []telegram.Request{
			entry.compareAsText(newer, destination),
		}
	case SingleMediaKind:
		return []telegram.Request{
			entry.compareAsSingleMedia(newer, destination),
		}
	}

	return
}

func (entry *ResultEntry) compareAgainst(
	newer Message,
	destination telegram.Destination,
	api *bot.Api,
	replaceMode bool,
) bool {

	defer func() {
		entry.Message = newer
	}()

	if replaceMode || entry.GetKind() != newer.GetKind() {
		replaceMode = true
		entry.cleanup(api, destination)
		response, err := api.Request(newer.InitRequest(destination))
		if err != nil {
			stacktraceBuf := make([]byte, 1000)
			runtime.Stack(stacktraceBuf, true)
			log.Printf("error in sending body message %v %v\n", err, stacktraceBuf)
			return replaceMode
		}

		messageList, _ := telegram.ParseAs[dto.MessageList](response)
		entry.BoundIds = messageList.CollectIds()
		return replaceMode
	}

	editRequests := entry.compareNonReplace(newer, destination)
	for _, request := range editRequests {
		api.LaunchRequest(request)
	}

	return replaceMode
}

func (result *Result) CompareAgainst(
	newer []Message,
	api *bot.Api,
	destination telegram.Destination,
	replaceMode bool,
) {
	commonLength := util.Min(len(result.Entries), len(newer))
	for i := 0; i < commonLength; i++ {
		replaceMode = result.Entries[i].compareAgainst(newer[i], destination, api, replaceMode)
	}

	if len(result.Entries) > len(newer) {
		for _, entry := range result.Entries[len(newer):] {
			entry.cleanup(api, destination)
		}
		result.Entries = result.Entries[:len(newer)]
	}

	if len(result.Entries) < len(newer) {
		//messageCount := len(newer) - len(result.Entries)

		result.Entries = append(result.Entries, util.Parallelize(func(newMessage Message) *ResultEntry {
			init := newMessage.InitRequest(destination)
			response, err := api.Request(init)
			if err != nil {
				api.LogError(err, init)
				return nil
			}
			messageList, _ := telegram.ParseAs[dto.MessageList](response)
			return &ResultEntry{
				Message:  newMessage,
				BoundIds: messageList.CollectIds(),
			}

		}, newer[len(result.Entries):])...)

	}

}