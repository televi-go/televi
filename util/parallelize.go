package util

func Parallelize[T any, TArg any](executor func(TArg) T, args []TArg) (result []T) {
	recChannel := make(chan T, len(args))
	for i := 0; i < len(args); i++ {
		go func(arg TArg) {
			recChannel <- executor(arg)
		}(args[i])
	}

	for i := 0; i < len(args); i++ {
		result = append(result, <-recChannel)
	}
	close(recChannel)

	return
}
