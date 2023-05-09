package util

func WhenBoth[T1 any, T2 any](firstRes T1, firstErr error, secondRes T2, secondErr error) (T1, T2, error) {
	if firstErr != nil {
		return firstRes, secondRes, firstErr
	}
	if secondErr != nil {
		return firstRes, secondRes, secondErr
	}
	return firstRes, secondRes, nil
}
