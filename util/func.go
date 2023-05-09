package util

func Map[T any, R any](source []T, transform func(elem T) R) []R {
	result := make([]R, len(source))

	for i, t := range source {
		result[i] = transform(t)
	}

	return result
}

func Filter[T any](source []T, filterFn func(elem T) bool) []T {
	result := make([]T, 0, len(source))
	for _, t := range source {
		if filterFn(t) {
			result = append(result, t)
		}
	}
	return result
}

func FilterOut[T any](source []T, filterFn func(elem T) bool) (conforming []T, other []T) {
	for _, t := range source {
		if filterFn(t) {
			conforming = append(conforming, t)
		} else {
			other = append(other, t)
		}
	}
	return
}

func MakePointerArr[T any](source []T) []*T {
	result := make([]*T, len(source))
	for i := 0; i < len(source); i++ {
		result[i] = &source[i]
	}
	return result
}

func UniqueEntries[T any, K comparable](source []T, picker func(T) K) []K {
	set := make(map[K]bool)

	for _, t := range source {
		picked := picker(t)
		set[picked] = true
	}

	result := make([]K, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	return result
}
