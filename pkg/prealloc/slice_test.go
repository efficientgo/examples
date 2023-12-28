// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package prealloc

import (
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestCreateSlice(t *testing.T) {
	const n = int(1e4)

	wp := createSlice(n)
	testutil.Equals(t, n*7, len(wp))

	// We can't predict exact capacity as it's based on many factors in Go runtime.
	testutil.Assert(t, n*7 <= cap(wp) && cap(wp) <= n*7+(8*1024))

	p := createSlice_Better(n)
	testutil.Equals(t, wp, p)
	testutil.Equals(t, n*7, len(p))
	testutil.Equals(t, len(p), cap(p))
}

func BenchmarkCreateSlice(b *testing.B) {
	const n = int(1e4)

	b.Run("without prealloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = createSlice(n)
		}
	})
	b.Run("prealloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = createSlice_Better(n)
		}
	})
}
