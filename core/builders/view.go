package builders

type View interface {
	Init()
	View(builder ComponentBuilder)
}

type DisposableView interface {
	View
	Dispose()
}

type ComponentBuilder interface {
	Component(View)
	Message(builder func(message Message))
}
