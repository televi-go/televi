package util

func Merge[T any](first, second []T) []T {
	result := make([]T, len(first)+len(second))
	for i := 0; i < len(first); i++ {
		result[i] = first[i]
	}
	for i := 0; i < len(second); i++ {
		result[i+len(first)] = second[i]
	}
	return result
}
