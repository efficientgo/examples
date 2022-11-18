package pools

import (
	"bytes"
	"io"
	"runtime"
	"sync"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestReuse(t *testing.T) {
	t.Run("", func(t *testing.T) {
		buf := make([]byte, 1e3)

		r := bytes.NewReader([]byte("abc"))
		w := &bytes.Buffer{}
		n, err := io.CopyBuffer(w, r, buf)

		testutil.Ok(t, err)
		testutil.Equals(t, 3, int(n))
		testutil.Equals(t, "abc", w.String())
	})
	t.Run("", func(t *testing.T) {
		buf := make([]byte, 1e3)

		r := bytes.NewReader([]byte("abc"))
		n, err := io.ReadFull(r, buf[:3])

		testutil.Ok(t, err)
		testutil.Equals(t, 3, int(n))
		testutil.Equals(t, "abc", string(buf[:n]))
	})
}

// BenchmarkProcess shows benchmarks of buffer that highlights common bug.
// Read more in "Efficient Go"; Example 11-20.
func BenchmarkProcess(b *testing.B) {
	b.Run("alloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			processUsingBuffer(nil)
		}
	})
	b.Run("buffer", func(b *testing.B) {
		b.ReportAllocs()
		buf := make([]byte, 1e6)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processUsingBuffer(buf)
		}
	})
	b.Run("pool-wrong", func(b *testing.B) {
		b.ReportAllocs()

		p := sync.Pool{
			New: func() any { return []byte{} },
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processUsingPool_Wrong(&p)
		}
	})
	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()

		p := sync.Pool{
			New: func() any { return []byte{} },
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processUsingPool(&p)
		}
	})
	b.Run("pool-GC", func(b *testing.B) {
		b.ReportAllocs()

		p := sync.Pool{
			New: func() any { return []byte{} },
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processUsingPool(&p)
			runtime.GC()
			runtime.GC()
		}
	})
	b.Run("buffer-GC", func(b *testing.B) {
		b.ReportAllocs()
		buf := make([]byte, 1e6)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processUsingBuffer(buf)
			runtime.GC()
			runtime.GC()
		}
	})
}
