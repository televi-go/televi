package test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"televi/models/pages"
	runner2 "televi/runner"
	"testing"
)

type MyPage struct {
	Counter   pages.State[int]
	Telephone pages.State[string]
}

func (m MyPage) View(ctx pages.PageBuildContext) {
	ctx.TextElement(func(component pages.TextContext) {
		component.Text(fmt.Sprintf("You clicked %d times", m.Counter.Get()))
		component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Increase", func(ctx pages.ReactionContext) {
					m.Counter.SetFn(func(prev int) int {
						return prev + 1
					})
				})
			})
		})
	})
	ctx.ActiveElement(func(component pages.ActiveTextContext) {
		component.TextF("UserInfo : %d\n", ctx.GetUserId()).Spoiler()
		if m.Telephone.Get() != "" {
			component.TextF("Telephone: %s", m.Telephone.Get())
		}
		component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ContactButton("Участвовать", func(ctx pages.ContactReactionContext) {
					m.Telephone.Set(ctx.Contact().PhoneNumber)
				})
			})
		})
	})
}

func TestRunner(t *testing.T) {
	runner, _ := runner2.NewRunner(os.Getenv("Token"), func() pages.Scene {
		return MyPage{}
	}, "root:@/televi?parseTime=true", runner2.DefaultAPiAddress)

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)

	defer cancel()

	runner.Run(ctx)
}
