package errorhandle

func Catch[T any, E any](r any, handler func(value T) (bool, E)) E {
	converted, ok := r.(T)
	if !ok {
		panic(r) // This should preserve the call stack as it is under the recover context
	}

	var result E
	isOK, result := handler(converted)
	if !isOK {
		panic(r)
	}

	return result
}

func Try[T any, U any, E any](functor func() U, handler func(value T) (bool, E)) (err E, result U) {
	defer func() {
		if r := recover(); r != nil {
			err = Catch(r, handler)
		}
	}()

	var nilErr E
	return nilErr, functor()
}
