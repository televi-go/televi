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

type MyFirstPage struct {
	Counter int
}

type MySecondPage struct {
}

func (secondPage *MySecondPage) View(ctx pages.PageBuildContext) {

	file, err := os.Open("photo_2023-04-03 23.01.17.jpeg")
	if err != nil {

		fmt.Println("error in opening file", err)
	}

	ctx.TextElement(func(component pages.TextContext) {
		component.Text("Some element")
	})
	ctx.PhotoElement(func(component pages.PhotoContext) {
		component.Image("my-photo", file)
		component.Text("Some other element")
		component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Forward", func(ctx pages.ReactionContext) {
					ctx.TransitTo(ThirdPage{}, pages.TransitPolicy{KeepPrevious: false})
				})
			})
		})
	})
}

type ThirdPage struct {
}

func (thirdPage ThirdPage) View(ctx pages.PageBuildContext) {
	ctx.TextElement(func(component pages.TextContext) {
		component.Text("This is third page only message")
		component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Transit back", func(ctx pages.ReactionContext) {
					ctx.TransitBack()
				})
			})
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Transit to main", func(ctx pages.ReactionContext) {
					ctx.TransitToMain()
				})
			})
		})
	})
}

func (firstPage *MyFirstPage) View(ctx pages.PageBuildContext) {
	ctx.TextElement(func(component pages.TextContext) {
		component.TextF("You clicked %d times", firstPage.Counter)
		component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Increase", func(ctx pages.ReactionContext) {
					firstPage.Counter++
				})
				rowBuilder.ActionButton("Transit", func(ctx pages.ReactionContext) {
					ctx.TransitTo(&MySecondPage{}, pages.TransitPolicy{KeepPrevious: true})
				})
			})
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Decrease", func(ctx pages.ReactionContext) {
					if firstPage.Counter <= 0 {
						ctx.ShowAlert(func(builder pages.TextPartBuilder) {
							builder.TextLine("Cannot remove")
							builder.TextLine("")
							builder.TextLine("Will be below zero")
						})
						return
					}
					firstPage.Counter--
				})
			})
		})
	})
}

func TestSeveralPages(t *testing.T) {
	runner, _ := runner2.NewRunner(os.Getenv("Token"), func() pages.Scene {
		return &MyFirstPage{}
	}, "root:@/televi?parseTime=true", runner2.DefaultAPiAddress)

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)

	defer cancel()

	runner.Run(ctx)
}
