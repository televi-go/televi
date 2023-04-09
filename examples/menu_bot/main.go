package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gtihub.com/televi-go/televi"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/runner"
	"log"
	"os"
)

type RootScene struct {
	// db field is needed for injecting into further scenes
	db *sql.DB
}

func (rootScene RootScene) View(ctx televi.BuildContext) {
	imageSource, _ := os.ReadFile("photo_2023-04-03 23.01.17.jpeg")
	ctx.ActivePhoto(func(component pages.ActivePhotoContext) {
		component.TextLine("Welcome to our progressive online restaurant")
		component.TextLine("This is our chef")
		component.Image("chef photo", bytes.NewReader(imageSource)).Spoiler()
		component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ActionButton("Menu", func(ctx pages.ReactionContext) {
					ctx.TransitTo(CategoriesScene{db: rootScene.db}, pages.TransitPolicy{KeepPrevious: false})
				})
				rowBuilder.ActionButton("Basket", func(ctx pages.ReactionContext) {
					ctx.TransitTo(BasketScene{service: basketService{db: rootScene.db}}, pages.TransitPolicy{KeepPrevious: false})
				})
			})
		})
	})
}

func main() {

	db, err := sql.Open("mysql", "root:@/televi?parseTime=true&multiStatements=true")

	if err != nil {
		log.Fatalln(err)
	}

	db.Exec(`
create temporary table categories(name varchar(32) not null primary key);
insert into categories (name) values('fish'),('bakery'),('drinks');
create temporary table basket(uid int not null, dish varchar(48) not null, primary key (uid, dish));
`)

	app, err := runner.NewRunner(os.Getenv("Token"), func() pages.Scene {
		return RootScene{db: db}
	}, "root:@/televi?parseTime=true", runner.DefaultAPiAddress)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run(context.TODO())
}
