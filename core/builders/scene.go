package builders

type Scene interface {
	Head(builder func(head Head))
	Body(builder func(body ComponentBuilder))
	//Navigator()
}
