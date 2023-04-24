package body

type Fragment struct {
	Messages []Message
}

type FragmentProducer interface {
	Build() Fragment
}
