package head_management

import (
	"github.com/televi-go/televi/core/callbacks"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages"
	"log"
)

type HeadResultContainer struct {
	Result        HeadResult
	HeadCallbacks callbacks.MenuCallbacks
	MessageId     int
}

func (container *HeadResultContainer) cleanup(destination telegram.Destination, api *bot.Api) {
	if container.MessageId != 0 {
		api.LaunchRequest(messages.DeleteMessageRequest{
			MessageId:   container.MessageId,
			Destination: destination,
		})
	}
}

func (container *HeadResultContainer) CompareAgainst(
	newer HeadResult,
	destination telegram.Destination,
	api *bot.Api,
) (wasReplaced bool) {
	if container.Result.equals(newer) {
		return false
	}
	container.cleanup(destination, api)

	request := newer.InitRequest(destination)
	response, err := api.Request(request)

	if err != nil {
		log.Printf("error painting head for %s: %v\n", destination.ToString(), err)
		return true
	}

	message, _ := telegram.ParseAs[dto.Message](response)
	container.MessageId = message.MessageID
	container.Result = newer
	return true
}
