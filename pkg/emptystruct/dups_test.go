package emptystruct

import (
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestHasDuplicates(t *testing.T) {
	testutil.Equals(t, true, HasDuplicates[float64](1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicates[float64](1, 2, 3, 4, 5, 6, 7))

	testutil.Equals(t, true, HasDuplicates_Better[float64](1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicates_Better[float64](1, 2, 3, 4, 5, 6, 7))

	testutil.Equals(t, true, HasDuplicates_NonGeneric(1, 2, 3, 4, 5, 6, 1))
	testutil.Equals(t, false, HasDuplicates_NonGeneric(1, 2, 3, 4, 5, 6, 7))
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
	b.Run("HasDuplicates_Better", func(b *testing.B) {
		b.ReportAllocs()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testutil.Equals(b, false, HasDuplicates_Better[float64](s...))
		}
	})
	b.Run("HasDuplicates_NonGeneric", func(b *testing.B) {
		b.ReportAllocs()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testutil.Equals(b, false, HasDuplicates_NonGeneric(s...))
		}
	})
}
