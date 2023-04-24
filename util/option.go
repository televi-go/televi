package util

type Option[T any] struct {
	value    T
	hasValue bool
}

func (option Option[T]) HasValue() bool {
	return option.hasValue
}

func (option Option[T]) ValuePtr() *T {
	if !option.hasValue {
		return nil
	}
	return &option.value
}

func OptionValue[T any](v T) Option[T] {
	return Option[T]{
		value:    v,
		hasValue: true,
	}
}

func TakeFirst[T any](source []T) Option[T] {
	if len(source) == 0 {
		return OptionEmpty[T]()
	}
	return OptionValue[T](source[0])
}

func OptionEmpty[T any]() Option[T] {
	return Option[T]{
		hasValue: false,
	}
}
