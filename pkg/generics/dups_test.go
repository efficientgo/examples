package generics

import (
	"testing"

	"github.com/efficientgo/core/testutil"
)

func HasDuplicates[T comparable](slice ...T) bool {
	dup := make(map[T]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}

func HasDuplicatesFloat64(slice ...float64) bool {
	dup := make(map[float64]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}

func TestHasDuplicates(t *testing.T) {
	testutil.Equals(t, true, HasDuplicates[float64](1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicates[float64](1, 2, 3, 4, 5, 6, 7))

	testutil.Equals(t, true, HasDuplicatesFloat64(1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicatesFloat64(1, 2, 3, 4, 5, 6, 7))
}

func BenchmarkHasDuplicates(b *testing.B) {
	s := make([]float64, 1e6)
	for i := range s {
		s[i] = 99.0 + float64(i)
	}

	b.Run("HasDuplicates", func(b *testing.B) {
		b.ReportAllocs()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testutil.Equals(b, false, HasDuplicates[float64](s...))
		}
	})
	b.Run("HasDuplicates2", func(b *testing.B) {
		b.ReportAllocs()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testutil.Equals(b, false, HasDuplicatesFloat64(s...))
		}
	})
}
