package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gtihub.com/televi-go/televi"
	"gtihub.com/televi-go/televi/models/pages"
)

type SelectDishScene struct {
	category     string
	service      basketService
	RenderHandle pages.State[int]
}

func (selectDishScene SelectDishScene) View(ctx televi.BuildContext) {
	allDishes, err := selectDishScene.service.getItems(ctx.GetUserId())
	userId := ctx.GetUserId()
	if err != nil {
		allDishes = map[string]bool{}
	}
	for i := 0; i < 5; i++ {
		i := i
		dishname := fmt.Sprintf("%s dish #%d", cases.Title(language.English).String(selectDishScene.category), i)
		ctx.TextElement(func(component pages.TextContext) {
			_, hasDish := allDishes[dishname]
			component.TextLine(dishname)
			if hasDish {
				component.TextLine("You have it in backet")
			}
			component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				if hasDish {
					rowBuilder.ActionButton("Remove", func(ctx pages.ReactionContext) {
						selectDishScene.service.toggleItem(userId, dishname)
						selectDishScene.RenderHandle.Set(0)
					})
					return
				}
				rowBuilder.ActionButton("Add", func(ctx pages.ReactionContext) {
					selectDishScene.service.toggleItem(userId, dishname)
					selectDishScene.RenderHandle.Set(0)
				})
			})

		})
	}
	ctx.TextElement(func(component pages.TextContext) {
		component.Text("Or return back to categories")
		component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
			rowBuilder.ActionButton("To categories", func(ctx pages.ReactionContext) {
				ctx.TransitBack()
			})
		})
	})
}
