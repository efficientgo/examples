package compileroptimizeaway

// Example of code that tends to be optimized by compiler in microbenchmarks.
// Thanks to Dave Cheney for this example: https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html#watch_out_for_compiler_optimisations
// Read more in "Efficient Go"; Example 8-16.

const m1 = 0x5555555555555555
const m2 = 0x3333333333333333
const m4 = 0x0f0f0f0f0f0f0f0f
const h01 = 0x0101010101010101

func popcnt(x uint64) uint64 {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return (x * h01) >> 56
}
