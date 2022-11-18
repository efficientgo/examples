package compileroptimizeaway

import (
	"math"
	"runtime"
	"testing"
)

// BenchmarkPopcnt_Wrong is an example microbenchmark that can be optimized by compiler.
// Read more in "Efficient Go"; Example 8-16.
func BenchmarkPopcnt_Wrong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcnt(math.MaxUint64)
	}
}

func BenchmarkPopcnt_Wrong2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcnt(Input)
	}
}

var Sink uint64

func BenchmarkPopcnt_Wrong3(b *testing.B) {
	var r uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = popcnt(math.MaxUint64)
	}
	Sink = r
}

// BenchmarkPopcnt_Sink is one example on how we can countermeasure the problem visible in BenchmarkPopcnt_Wrong.
// Read more in "Efficient Go"; Example 8-18.
func BenchmarkPopcnt_Sink(b *testing.B) {
	var r uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = popcnt(Input)
	}
	Sink = r
}

func BenchmarkPopcnt_KeepAlive(b *testing.B) {
	var r uint64

	for i := 0; i < b.N; i++ {
		r = popcnt(Input)
	}
	runtime.KeepAlive(r)
}

var Input uint64 = math.MaxUint64
