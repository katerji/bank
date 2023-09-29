package utils

func InSlice[T comparable](s []T, itemToSearchFor T) bool {
	for _, item := range s {
		if itemToSearchFor == item {
			return true
		}
	}
	return false
}
