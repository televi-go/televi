package main

import (
	"database/sql"
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/models/pages"
)

type BasketScene struct {
	service basketService
}

func (basketScene BasketScene) View(ctx televi.BuildContext) {
	items, err := basketScene.service.getItems(ctx.GetUserId())

	ctx.TextElement(func(component pages.TextContext) {
		if err != nil {
			component.Text("Error encountered")
			component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
				rowBuilder.ActionButton("Go back", func() {
					ctx.GetNavigator().Pop()
				})
			})
			return
		}
		component.TextLine("Your basket:").Bold()
		for dish, _ := range items {
			component.TextF("%s\n", dish)
		}
		component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
			rowBuilder.ActionButton("Checkout", func() {
				basketScene.service.clear(ctx.GetUserId())
				ctx.GetNavigator().Alert(func(builder pages.TextPartBuilder) {
					builder.TextLine("Your purchase list:\n")
					for dish, _ := range items {
						builder.TextF("%s\n", dish)
					}
					builder.TextLine("Enjoy your meal!")
				})
				ctx.GetNavigator().PopAll()
			})
		})
	})
}

type basketService struct {
	db       *sql.DB
	onUpdate func()
}

func (service basketService) getItems(uid int) (map[string]bool, error) {
	rows, err := service.db.Query("select dish from basket where uid=?", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]bool)

	for rows.Next() {
		var dish string
		err = rows.Scan(&dish)
		if err != nil {
			return nil, err
		}
		result[dish] = true
	}

	return result, nil
}

func (service basketService) toggleItem(uid int, item string) (err error) {
	defer func() {
		if err == nil {
			service.onUpdate()
		}
	}()
	existRow := service.db.QueryRow("select exists(select *from basket where dish=? and uid=?)", item, uid)
	var exists bool
	err = existRow.Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = service.db.Exec("delete from basket where dish=? and uid=?", item, uid)
		return err
	}

	_, err = service.db.Exec("insert into basket (uid, dish) values (?, ?)", uid, item)

	return err
}

func (service basketService) clear(uid int) error {
	_, err := service.db.Exec("delete from basket where uid=?", uid)
	return err
}
