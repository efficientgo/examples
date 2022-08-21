package pools

import (
	"bytes"
	"io"
	"runtime"
	"sync"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
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

func processUsingBuffer(buf []byte) {
	buf = buf[:0]

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}

	// Use buffer...
}

func processUsingPool_Wrong(p *sync.Pool) {
	buf := p.Get().([]byte)
	buf = buf[:0]

	defer p.Put(buf)

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}

	// Use buffer...
}

func processUsingPool(p *sync.Pool) {
	buf := p.Get().([]byte)
	buf = buf[:0]

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}
	defer p.Put(buf)

	// Use buffer...
}

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
