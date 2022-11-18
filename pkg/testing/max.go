package max

// max returns maximum of a and b.
// If both are equal, returns a.
//
// Used as an example for unit tests.
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
