package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/runner"
	"log"
	"os"
	"os/signal"
)

type RootScene struct {
	// db field is needed for injecting into further scenes
	db *sql.DB
}

func (rootScene RootScene) View(ctx televi.BuildContext) {

	ctx.ActivePhoto(func(component pages.ActivePhotoContext) {
		component.TextF("Hi, %s\n", ctx.GetUserInfo().FirstName)
		component.TextLine("Welcome to our progressive online restaurant")
		component.TextLine("This is our chef")
		chefAsset.Embed(component).Spoiler()
		component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ActionButton("Menu", func() {
					ctx.GetNavigator().Extend(CategoriesScene{db: rootScene.db})
				})
				rowBuilder.ActionButton("Basket", func() {
					ctx.GetNavigator().Extend(BasketScene{
						service: basketService{
							db: rootScene.db,
						},
					})
				})
			})
		})
	})
}

var chefAsset televi.ImageAsset
var dishAssets = map[string]*[5]televi.ImageAsset{
	"fish":   {},
	"drinks": {},
	"bakery": {},
}

var dishAssetsPath = map[string][5]string{
	"bakery": {"IMG_4639 2.jpg", "IMG_4640 2.jpg", "IMG_4641 2.jpg", "IMG_4642 3.jpg", "IMG_4642 4.jpg"},
	"drinks": {"IMG_4645.jpg", "IMG_4645 2.jpg", "IMG_4645 3.jpg", "IMG_4646.jpg", "IMG_4649.jpg"},
	"fish":   {"IMG_4650.jpg", "IMG_4650 2.jpg", "IMG_4650 3.jpg", "IMG_4650 4.jpg", "IMG_4651.jpg"},
}

func main() {

	db, err := sql.Open("mysql", "root:@/televi?parseTime=true&multiStatements=true")

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	db.Exec(`
create temporary table categories(name varchar(32) not null primary key);
insert into categories (name) values('fish'),('bakery'),('drinks');
create temporary table basket(uid int not null, dish varchar(48) not null, primary key (uid, dish));
`)

	assetLoader := televi.NewAssetLoader()

	assetLoader.Add("examples/menu_bot/assets/photo_2023-04-03 23.01.17.jpeg", &chefAsset)

	for category, assets := range dishAssets {
		for i := 0; i < 5; i++ {
			assetLoader.Add("examples/menu_bot/assets/"+dishAssetsPath[category][i], &(*assets)[i])
		}
	}

	err = assetLoader.Load()

	if err != nil {
		log.Fatalln("Error loading asset", err)
	}

	app, err := runner.NewRunner(os.Getenv("Token"), func() pages.Scene {
		return RootScene{db: db}
	}, "root:@/televi?parseTime=true", runner.EnvOrDefault("Address"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("starting app")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	app.Run(ctx)
}
