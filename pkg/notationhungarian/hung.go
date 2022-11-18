package notationhungarian

// structSystem represents old way of naming code structures (putting type in the name). It's anti-pattern nowadays.
// Read more in "Efficient Go"; Example 1-5.
type structSystem struct {
	sliceU32Numbers []uint32
	bCharacter      byte
	f64Ratio        float64
}
