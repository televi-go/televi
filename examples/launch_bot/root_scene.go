package main

import (
	"database/sql"
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/models/pages"
)

type RootScene struct {
	db *sql.DB
}

func (rootScene RootScene) View(ctx televi.BuildContext) {
	assetLoader := televi.NewAssetLoader()
	assetLoader.
		Add("examples/launch_bot/welcome_pic.png", &welcomePicAsset).
		Add("examples/launch_bot/animation_sub.gif", &subScribedAsset)
	_ = assetLoader.Load()
	ctx.ActivePhoto(func(component pages.ActivePhotoContext) {

		welcomePicAsset.Embed(component)

		component.TextF("Hi, %s\n", ctx.GetUserInfo().FirstName)
		component.TextLine("Welcome to our study course")
		component.TextLine("Hit the button to subscribe")

		component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ActionButton("Subscribe to webinar", func() {
					ctx.GetNavigator().Replace(SubscribedScene{})
				})
			})
		})

	})

}

var welcomePicAsset televi.ImageAsset
var subScribedAsset televi.ImageAsset
