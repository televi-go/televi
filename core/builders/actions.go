package builders

type ActionsBuilder interface {
	Row(func(builder ActionRowBuilder))
}

type ActionRowBuilder interface {
	Button(caption string, onclick func())
	Url(caption string, target string)
}
