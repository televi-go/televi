package load_testing

import (
	"github.com/televi-go/televi/runner"
	"github.com/televi-go/televi/telegram/dto"
	"golang.org/x/exp/rand"
	"sync"
)

func SyntheticLoadRunner(runner *runner.Runner, count int) {
	wg := sync.WaitGroup{}

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			runner.DispatchUpdate(dto.Update{
				UpdateID: 0,
				Message: &dto.Message{
					MessageID: rand.Int(),
					From: &dto.User{
						ID:                      int64(rand.Uint64()),
						IsBot:                   false,
						IsPremium:               false,
						FirstName:               "",
						LastName:                "",
						UserName:                "",
						LanguageCode:            "",
						CanJoinGroups:           false,
						CanReadAllGroupMessages: false,
						SupportsInlineQueries:   false,
					},
					SenderChat: nil,
					Date:       0,
					Chat: &dto.Chat{
						ID:                    int64(rand.Int()),
						Type:                  "",
						Title:                 "",
						UserName:              "",
						FirstName:             "",
						LastName:              "",
						Photo:                 nil,
						Bio:                   "",
						HasPrivateForwards:    false,
						Description:           "",
						InviteLink:            "",
						PinnedMessage:         nil,
						Permissions:           nil,
						SlowModeDelay:         0,
						MessageAutoDeleteTime: 0,
						HasProtectedContent:   false,
						StickerSetName:        "",
						CanSetStickerSet:      false,
						LinkedChatID:          0,
						Location:              nil,
					},
					ForwardFrom:             nil,
					ForwardFromChat:         nil,
					ForwardFromMessageID:    0,
					ForwardSignature:        "",
					ForwardSenderName:       "",
					ForwardDate:             0,
					IsAutomaticForward:      false,
					ReplyToMessage:          nil,
					ViaBot:                  nil,
					EditDate:                0,
					HasProtectedContent:     false,
					MediaGroupID:            "",
					AuthorSignature:         "",
					Text:                    "abcefg",
					Entities:                nil,
					Animation:               nil,
					PremiumAnimation:        nil,
					Audio:                   nil,
					Document:                nil,
					Photo:                   nil,
					Sticker:                 nil,
					Video:                   nil,
					VideoNote:               nil,
					Voice:                   nil,
					Caption:                 "",
					CaptionEntities:         nil,
					Contact:                 nil,
					Dice:                    nil,
					Poll:                    nil,
					Venue:                   nil,
					Location:                nil,
					NewChatMembers:          nil,
					LeftChatMember:          nil,
					NewChatTitle:            "",
					NewChatPhoto:            nil,
					DeleteChatPhoto:         false,
					GroupChatCreated:        false,
					SuperGroupChatCreated:   false,
					ChannelChatCreated:      false,
					MigrateToChatID:         0,
					MigrateFromChatID:       0,
					PinnedMessage:           nil,
					ConnectedWebsite:        "",
					ProximityAlertTriggered: nil,
				},
				EditedMessage:      nil,
				ChannelPost:        nil,
				EditedChannelPost:  nil,
				InlineQuery:        nil,
				ChosenInlineResult: nil,
				CallbackQuery:      nil,
				ShippingQuery:      nil,
				PreCheckoutQuery:   nil,
				Poll:               nil,
				PollAnswer:         nil,
				MyChatMember:       nil,
				ChatMember:         nil,
				ChatJoinRequest:    nil,
			})
		}()
	}
	wg.Wait()

}
