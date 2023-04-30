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
