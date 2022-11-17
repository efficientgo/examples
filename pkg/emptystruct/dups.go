package emptystruct

// Example of simple optimization that allows to do less work which is not necessary.
// Read more in "Efficient Go"; Example 11-1.

func HasDuplicates[T comparable](slice ...T) bool {
	dup := make(map[T]any, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = "whatever, I don't use this value"
	}
	return false
}

func HasDuplicates_Better[T comparable](slice ...T) bool {
	dup := make(map[T]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}

func HasDuplicates_NonGeneric(slice ...float64) bool {
	dup := make(map[float64]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}
