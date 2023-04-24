package main

import (
	"github.com/televi-go/televi"
	"github.com/televi-go/televi/models/pages"
	"time"
)

type SubscribedScene struct {
	ShowAnnouncement pages.State[bool]
}

func (subScene SubscribedScene) Init() {

	// emulate any deadline
	<-time.After(time.Second * 10)
	subScene.ShowAnnouncement.Set(true)

}

type InterestedScene struct {
	cost     int
	tariff   string
	TimeLeft pages.State[time.Duration]
}

func (interestedScene InterestedScene) Init() {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for range timer.C {
		if interestedScene.TimeLeft.Get() == 0 {
			return
		}

		interestedScene.TimeLeft.SetFn(func(prev time.Duration) time.Duration {
			return prev - time.Second
		})
	}
}

func (interestedScene InterestedScene) View(ctx televi.BuildContext) {
	ctx.TextElement(func(component pages.TextContext) {

		var price = interestedScene.cost

		if interestedScene.TimeLeft.Get() > 0 {
			price = price * 70 / 100
		}

		component.TextLine("Alright, you have shown your interest.")
		dur := time.Time{}.Add(interestedScene.TimeLeft.Get())
		component.TextF("Discount will end in %s\n", dur.Format(time.TimeOnly))
		component.TextF("Your price: %d", price).Bold()

	})
	ctx.ActiveElement(func(component pages.ActiveTextContext) {
		component.TextLine("You'd better buy it now")
		component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ActionButton("Buy", func() {

				})
			})
			builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
				rowBuilder.ActionButton("Select another tariff", func() {
					ctx.GetNavigator().Pop()
				})
			})
		})
	})
}

type AnnouncementScene struct {
	WebinarHeld pages.State[bool]
}

func (announcementScene AnnouncementScene) Init() {
	<-time.After(time.Second * 5)
	announcementScene.WebinarHeld.Set(true)
}

func (announcementScene AnnouncementScene) View(ctx televi.BuildContext) {
	ctx.TextElement(func(component pages.TextContext) {
		component.TextLine("Hey, we are holding a webinar")
		component.ButtonsRow(func(rowBuilder pages.InlineRowBuilder) {
			rowBuilder.UrlButton("Here it will be", "https://google.com")
		})
	})
	if announcementScene.WebinarHeld.Get() {
		ctx.ActiveElement(func(component pages.ActiveTextContext) {
			component.TextLine("How did you like our webinar?")
			component.TextLine("We are sure you are going for more")
			component.TextLine("Here are three options")
			component.ReplyKeyboard(func(builder pages.ReplyKeyboardBuilder) {
				builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
					rowBuilder.ActionButton("Common", func() {
						ctx.GetNavigator().Extend(InterestedScene{TimeLeft: pages.StateOf(time.Second * 100), cost: 1000})
					})
				})
				builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
					rowBuilder.ActionButton("Lite", func() {
						ctx.GetNavigator().Extend(InterestedScene{TimeLeft: pages.StateOf(time.Second * 100), cost: 100})
					})
				})
				builder.ButtonsRow(func(rowBuilder pages.ReplyRowBuilder) {
					rowBuilder.ActionButton("Premium", func() {
						ctx.GetNavigator().Extend(InterestedScene{TimeLeft: pages.StateOf(time.Second * 100), cost: 3000})
					})
				})
			})
		})
	}
}

func (subScene SubscribedScene) View(ctx televi.BuildContext) {
	ctx.ActiveElement(func(component pages.ActiveTextContext) {
		component.TextLine("You have subscribed!\n").Bold()
	})
	if subScene.ShowAnnouncement.Get() {
		ctx.GetNavigator().Replace(AnnouncementScene{})
	}
}
