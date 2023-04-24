package magic

import "unsafe"

type State[T any] struct {
	Trigger
	value *T
}

func (s State[T]) SetValueFn(changer func(previous T) T) {
	s.SetValue(changer(s.Value()))
}

func (s State[T]) Value() T {
	return *s.value
}

func (s State[T]) SetValue(value T) {
	*s.value = value
	s.SendChange()
}

func (t Trigger) SendChange() {
	t.onChange()
}

func (s State[T]) Mount(c func(), mountPtr unsafe.Pointer) {
	ptr := (*State[T])(mountPtr)
	ptr.onChange = c
	if ptr.value == nil {
		ptr.value = new(T)
	}
}

func (t Trigger) Mount(c func(), mountPtr unsafe.Pointer) {
	(*Trigger)(mountPtr).onChange = c
}

type Trigger struct {
	onChange func()
}

type Mountable interface {
	Mount(func(), unsafe.Pointer)
}
