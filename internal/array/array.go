package array

func Prepend[T any](arr []T, val T) []T {
	return append([]T{val}, arr...)
}

func Contains[T comparable](arr []T, val T) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
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
