package telegram

import (
	"github.com/televi-go/televi/telegram/dto"
	"strconv"
)

// Destination is a marker interface representing communication end
type Destination interface {
	ParamsWriter
	ToString() string
	destinationImpl()
}

type ChatDestination struct {
	ChatId int
}

func (chatDestination ChatDestination) ToString() string {
	return strconv.Itoa(chatDestination.ChatId)
}

func (chatDestination ChatDestination) destinationImpl() {}

type ChannelDestination struct {
	ChannelName string
}

func (channelDestination ChannelDestination) ToString() string {
	return channelDestination.ChannelName
}

func (channelDestination ChannelDestination) destinationImpl() {}

func (channelDestination ChannelDestination) WriteParameter(params Params) error {
	params.WriteString("chat_id", channelDestination.ChannelName)
	return nil
}

func (chatDestination ChatDestination) WriteParameter(params Params) error {
	params.WriteInt("chat_id", chatDestination.ChatId)
	return nil
}

func GetDestination(update dto.Update) Destination {
	if update.Message != nil {
		return ChatDestination{ChatId: int(update.Message.Chat.ID)}
	}

	if update.CallbackQuery != nil {
		return ChatDestination{ChatId: int(update.CallbackQuery.From.ID)}
	}

	panic("cannot resolve destination")
}

func ParseDestination(source string) Destination {
	chat, err := strconv.Atoi(source)
	if err != nil {
		return ChannelDestination{ChannelName: source}
	}
	return ChatDestination{ChatId: chat}
}
