package builders

type Scene interface {
	Head(builder func(headBuilder Head))
	Body(builder func(bodyBuilder ComponentBuilder))
	Navigator()
}
