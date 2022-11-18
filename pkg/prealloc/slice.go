package prealloc

// createSlice represents less efficient (also less readable IMO) creation of slice.
// Read more in "Efficient Go"; Example 1-3.
func createSlice(n int) (slice []string) {
	for i := 0; i < n; i++ {
		slice = append(slice, "I", "am", "going", "to", "take", "some", "space")
	}
	return slice
}

// createSlice_Better represents more efficient creation of slice that pre-allocates array in memory upfront.
// Read more in "Efficient Go"; Example 1-4.
func createSlice_Better(n int) []string {
	slice := make([]string, 0, n*7)
	for i := 0; i < n; i++ {
		slice = append(slice, "I", "am", "going", "to", "take", "some", "space")
	}
	return slice
}
