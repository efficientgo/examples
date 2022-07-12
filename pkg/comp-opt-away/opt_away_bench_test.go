package comp_opt_away

import (
	"math"
	"runtime"
	"testing"
)

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

func BenchmarkPopcnt_KeepAlive(b *testing.B) {
	var r uint64

	for i := 0; i < b.N; i++ {
		r = popcnt(Input)
	}
	runtime.KeepAlive(r)
}

var Input uint64 = math.MaxUint64

func BenchmarkPopcnt_Sink(b *testing.B) {
	var r uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = popcnt(Input)
	}
	Sink = r
}
