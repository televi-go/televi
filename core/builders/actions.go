package builders

type ActionsBuilder interface {
	Row(func(builder ActionRowBuilder))
}

type DisplayMode int

const (
	AlertMode        DisplayMode = iota
	NotificationMode DisplayMode = iota
)

type AlertBuilder interface {
	ContentBuilder
	SetDisplayMode(mode DisplayMode)
}

type ClickContext interface {
	Alert(builder func(AlertBuilder))
}

type ActionRowBuilder interface {
	Button(caption string, onclick func(context ClickContext))
	Url(caption string, target string)
}
