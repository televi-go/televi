package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/load_testing"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/runner"
	"log"
	"os"
	"os/signal"
	"testing"
)

func TestBotLoad(t *testing.T) {
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
	go load_testing.SyntheticLoadRunner(app, 1_000)
	app.Run(ctx)

}
