package utils

// Convert slice to slice of pointers
func ToPointers[T any](arr []T) []*T {
	result := make([]*T, len(arr))
	for a := range len(arr) {
		result[a] = &arr[a]
	}

	return result
}

func Map[T, R any](arr []T, mapper func(T) R) []R {
	result := make([]R, len(arr))
	for index, element := range arr {
		result[index] = mapper(element)
	}
	return result
}
