package utils

type Predictor[T any] func(val T) bool

func Filter[T any](slice []T, predictor Predictor[T]) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predictor(v) {
			result = append(result, v)
		}
	}

	return result
}

func FilterNotNil[T any](slice []T) []T {
	return Filter(slice, func(val T) bool { return !IsNil(val) })
}

func Any[T any](slice []T, predictor Predictor[T]) bool {
	for _, v := range slice {
		if predictor(v) {
			return true
		}
	}
	return false
}

func None[T any](slice []T, predictor Predictor[T]) bool {
	for _, v := range slice {
		if predictor(v) {
			return false
		}
	}
	return true
}
