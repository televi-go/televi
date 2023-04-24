package builders

type View interface {
	Init()
	View(builder ComponentBuilder)
}

type ComponentBuilder interface {
	Component(View)
	Message(builder func(viewBuilder Message))
}
