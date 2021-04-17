package prealloc

// Example showing readability difference between slice creation with and without pre-allocation.
// You can read more in "Efficient Go" book, Chapter 1.

func CreateSlice(n int) (slice []string) {
	for i := 0; i < n; i++ {
		slice = append(slice, "I", "am", "going", "to", "take", "some", "space")
	}
	return slice
}

func CreateSlice2(n int) []string {
	slice := make([]string, 0, n*7)
	for i := 0; i < n; i++ {
		slice = append(slice, "I", "am", "going", "to", "take", "some", "space")
	}
	return slice
}
