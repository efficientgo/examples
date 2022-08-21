package emptystruct

import (
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

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

func HasDuplicates2[T comparable](slice ...T) bool {
	dup := make(map[T]struct{}, len(slice))
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

	testutil.Equals(t, true, HasDuplicates2[float64](1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicates2[float64](1, 2, 3, 4, 5, 6, 7))
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
			testutil.Equals(b, false, HasDuplicates2[float64](s...))
		}
	})
}
