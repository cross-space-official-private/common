package utils

type Transformer[T any, P any] func(val T) P

func Map[T any, P any](slice []T, transformer Transformer[T, P]) []P {
	result := make([]P, len(slice))
	for i, v := range slice {
		result[i] = transformer(v)
	}

	return result
}

func FlatMap[T any, P any](slice [][]T, transformer Transformer[T, P]) []P {
	result := make([]P, 0)
	for _, v := range slice {
		result = append(result, Map(v, transformer)...)
	}

	return result
}

func Contains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateBy[T any, P comparable](sliceList []T, transformer Transformer[T, P]) []T {
	allKeys := make(map[P]bool)
	var list []T
	for _, item := range sliceList {
		key := transformer(item)
		if _, value := allKeys[key]; !value {
			allKeys[key] = true
			list = append(list, item)
		}
	}
	return list
}
