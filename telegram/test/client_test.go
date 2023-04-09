package test

import (
	"context"
	"fmt"
	"os"
	"televi/runner"
	"televi/telegram"
	"televi/telegram/bot"
	"televi/telegram/dto"
	"televi/telegram/messages"
	"testing"
	"time"
	"unsafe"
)

func TestGetUpdates(t *testing.T) {
	api := bot.NewApi(os.Getenv("Token"), runner.DefaultAPiAddress)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	for update := range api.Poll(ctx) {
		if update.Message != nil {
			fmt.Printf("Send message to %#v\n", update.Message.From)
			resp, err := api.Request(messages.TextMessageRequest{
				Destination: telegram.ChatDestination{
					ChatId: int(update.Message.Chat.ID),
				},
				Text:        "AbCD",
				Silent:      false,
				ReplyMarkup: nil,
			})
			fmt.Println(resp, err)
		}
	}
}

func TestMessagesSize(t *testing.T) {
	message := dto.Message{}
	fmt.Println(unsafe.Sizeof(message))
}
