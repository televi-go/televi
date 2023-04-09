package main

import (
	"context"
	"log"
	"os"
	"televi/models/pages"
	runner "televi/runner"
)

type SamplePage struct {
	Counter int
}

func (samplePage *SamplePage) View(ctx pages.PageBuildContext) {
	ctx.TextElement(func(component pages.TextContext) {
		component.TextF("Your value is %d", samplePage.Counter)
		component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Increase", func(ctx pages.ReactionContext) {
					samplePage.Counter++
				})
			})
		})
	})
}

func main() {
	app, err := runner.NewRunner(os.Getenv("Token"), func() pages.Scene {
		return &SamplePage{}
	}, "root:@/televi?parseTime=true", "http://localhost:8081")
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(context.TODO())
}
