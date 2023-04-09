package main

import (
	"database/sql"
	"fmt"
	"gtihub.com/televi-go/televi"
	"gtihub.com/televi-go/televi/models/pages"
)

type CategoriesScene struct {
	db            *sql.DB
	ReloadManager pages.State[struct{}]
}

func (categoriesScene CategoriesScene) getCategories() ([]string, error) {
	rows, err := categoriesScene.db.Query("select name from categories")
	if err != nil {
		return nil, err
	}
	var result []string
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		result = append(result, category)
	}
	return result, nil
}

func (categoriesScene CategoriesScene) View(ctx televi.BuildContext) {
	categories, err := categoriesScene.getCategories()
	if err != nil {
		fmt.Println("User", ctx.GetUserId(), "encountered", err)
		ctx.TextElement(func(component pages.TextContext) {
			component.Text("Service is unavailable. Try again later")
			component.InlineKeyboard(func(builder pages.InlineKeyboardBuilder) {
				builder.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
					rowBuilder.ActionButton("Reload", func(ctx pages.ReactionContext) {
						categoriesScene.ReloadManager.Set(struct{}{})
					})
				})
			})
		})
		return
	}
	ctx.TextElement(func(component pages.TextContext) {
		component.Text("Pick a category")
		for _, category := range categories {
			category := category
			component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton(category, func(ctx pages.ReactionContext) {
					ctx.TransitTo(SelectDishScene{
						category: category,
						service: basketService{
							db: categoriesScene.db,
						},
					}, pages.TransitPolicy{KeepPrevious: true})
				})
			})
		}
	})
}
