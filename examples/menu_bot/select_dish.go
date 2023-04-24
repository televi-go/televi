package main

import (
	"fmt"
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/models/pages"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SelectDishScene struct {
	category     string
	service      *basketService
	RenderHandle pages.State[int]
}

func (selectDishScene SelectDishScene) bindService() {
	selectDishScene.service.onUpdate = func() {
		selectDishScene.RenderHandle.Set(0)
	}
}

type DishView struct {
	hasDish  bool
	category string
	index    int
	dishName string
	service  basketService
}

func (dishView DishView) Build(ctx pages.PageBuildContext) {
	ctx.PhotoElement(func(component pages.PhotoContext) {
		dishAssets[dishView.category][dishView.index].Embed(component)
		component.TextLine(dishView.dishName)
		if dishView.hasDish {
			component.TextLine("You have it in basket")
		}
		component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
			if dishView.hasDish {
				rowBuilder.ActionButton("Remove", func() {
					dishView.service.toggleItem(ctx.GetUserId(), dishView.dishName)
				})
				return
			}
			rowBuilder.ActionButton("Add", func() {
				dishView.service.toggleItem(ctx.GetUserId(), dishView.dishName)
			})
		})
	})
}

func (selectDishScene SelectDishScene) View(ctx televi.BuildContext) {
	allDishes, err := selectDishScene.service.getItems(ctx.GetUserId())
	selectDishScene.bindService()
	if err != nil {
		allDishes = map[string]bool{}
	}
	televi.ForEach(ctx, televi.Range(0, 5), func(i int) pages.View {
		dishname := fmt.Sprintf("%s dish #%d", cases.Title(language.English).String(selectDishScene.category), i)
		_, hasDish := allDishes[dishname]
		return DishView{
			hasDish:  hasDish,
			category: selectDishScene.category,
			index:    i,
			dishName: dishname,
			service:  *selectDishScene.service,
		}
	})
	ctx.TextElement(func(component pages.TextContext) {
		component.Text("Or return back to categories")
		component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
			rowBuilder.ActionButton("To categories", func() {
				ctx.GetNavigator().Pop()
			})
		})
	})
}
