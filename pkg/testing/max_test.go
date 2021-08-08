package max

import (
	"math"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestMax(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		a, b     int
		expected int
	}{
		{a: 0, b: 0, expected: 0},
		{a: -1, b: 0, expected: 0},
		{a: 1, b: 0, expected: 1},
		{a: 0, b: -1, expected: 0},
		{a: 0, b: 1, expected: 1},
		{a: math.MinInt64, b: math.MaxInt64, expected: math.MaxInt64},
	} {
		t.Run("", func(t *testing.T) {
			testutil.Equals(t, tcase.expected, max(tcase.a, tcase.b))
		})
	}
}

// TODO(bwplotka): Explain when fuzzing will be part of Go.
// See: https://github.com/golang/go/issues/44551#issuecomment-811607377
//func FuzzMarshalFoo(f *testing.F) {
//
//}
