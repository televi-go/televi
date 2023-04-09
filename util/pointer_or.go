package util

import "fmt"

func PointerOr[T any](v any) (*T, error) {
	ptr, isPtr := v.(*T)
	if isPtr {
		return ptr, nil
	}
	val, isVal := v.(T)
	if isVal {
		return &val, nil
	}
	var t *T
	return nil, fmt.Errorf("value %+v does not conform to type %T", v, t)
}
