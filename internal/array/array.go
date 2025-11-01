package array

func Prepend[T any](arr []T, val T) []T {
	return append([]T{val}, arr...)
}

func ContainsMoreThan[T comparable](arr []T, values ...T) bool {
	for _, item := range arr {
		success := true
		for _, val := range values {
			if item == val {
				success = false
			}
		}
		if success {
			return true
		}
	}
	return false
}
