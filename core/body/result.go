package body

import (
	"bytes"
	"github.com/televi-go/televi/models/render/results"
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
	Interruption Interruption
	BoundIds     []int
}

type Result struct {
	Entries             []*ResultEntry
	PendingInterruption Interruption
	BoundInterruption   Interruption
}

func (result *Result) AddInterruption(id int, isNav bool) {
	result.PendingInterruption.MessageIds = append(result.PendingInterruption.MessageIds, InterEntry{MessageId: id, IsNavigation: isNav})
}

type InterEntry struct {
	MessageId    int
	IsNavigation bool
}

type Interruption struct {
	MessageIds []InterEntry
}

func (i Interruption) IsEmpty() bool {
	return len(i.MessageIds) == 0
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
	for _, id := range entry.Interruption.MessageIds {
		if id.MessageId != 0 {
			api.LaunchRequest(messages.DeleteMessageRequest{MessageId: id.MessageId, Destination: destination})
		}
	}
	entry.Interruption = Interruption{}
}

func areEqualActionKeyboards(first, second results.InlineKeyboardResult) bool {
	if len(first.Keyboard) != len(second.Keyboard) {
		return false
	}
	for i := 0; i < len(first.Keyboard); i++ {
		if len(first.Keyboard[i]) != len(second.Keyboard[i]) {
			return false
		}
		for j := 0; j < len(first.Keyboard[i]); j++ {
			if first.Keyboard[i][j] != second.Keyboard[i][j] {
				return false
			}
		}
	}
	return true
}

func (entry *ResultEntry) compareAsText(
	newer Message, destination telegram.Destination) util.Option[telegram.Request] {

	if entry.Text == newer.Text && areEqualActionKeyboards(entry.Actions, newer.Actions) {
		return util.OptionEmpty[telegram.Request]()
	}

	return util.OptionValue[telegram.Request](messages.EditMessageRequest{
		Destination: destination,
		Text:        newer.Text,
		MessageId:   entry.BoundIds[0],
		ReplyMarkup: newer.Actions.ToReplyMarkup(),
	})
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

	if areEqualActionKeyboards(entry.Actions, newer.Actions) {
		return nil
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
		return entry.compareAsText(newer, destination).ToSlice()
	case SingleMediaKind:
		return util.OptionFromNullable[telegram.Request](entry.compareAsSingleMedia(newer, destination)).ToSlice()
	}

	return
}

func (entry *ResultEntry) compareReplacing(
	newer Message,
	destination telegram.Destination,
	api *bot.Api,
) {
	entry.cleanup(api, destination)
	response, err := api.Request(newer.InitRequest(destination))
	if err != nil {
		stacktraceBuf := make([]byte, 1000)
		runtime.Stack(stacktraceBuf, true)
		log.Printf("error in sending body message %v %v\n", err, stacktraceBuf)
		return
	}

	messageList, _ := telegram.ParseAs[dto.MessageList](response)
	entry.BoundIds = messageList.CollectIds()
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
		entry.compareReplacing(newer, destination, api)
		return replaceMode
	}

	editRequests := entry.compareNonReplace(newer, destination)

	if len(editRequests) != 0 && len(entry.Interruption.MessageIds) != 0 {
		entry.compareReplacing(newer, destination, api)
		entry.Interruption = Interruption{}
		return true
	}

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

	if !result.BoundInterruption.IsEmpty() && !result.PendingInterruption.IsEmpty() {
		replaceMode = true
		for _, id := range result.BoundInterruption.MessageIds {
			api.LaunchRequest(messages.DeleteMessageRequest{MessageId: id.MessageId, Destination: destination})
		}
		result.BoundInterruption = Interruption{}
	}

	if len(result.PendingInterruption.MessageIds) != 0 {
		//replaceMode = true

		for _, interEntry := range result.PendingInterruption.MessageIds {
			if interEntry.IsNavigation || len(result.Entries) == 0 {
				result.BoundInterruption.MessageIds = append(result.BoundInterruption.MessageIds, interEntry)
			} else {
				result.Entries[len(result.Entries)-1].Interruption.MessageIds =
					append(result.Entries[len(result.Entries)-1].Interruption.MessageIds, interEntry)
			}
		}

		result.PendingInterruption = Interruption{}
	}

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
