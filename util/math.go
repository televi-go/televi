package util

import "golang.org/x/exp/constraints"

func Min[T constraints.Integer | constraints.Float](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Integer | constraints.Float](x, y T) T {
	if x > y {
		return x
	}
	return y
}
