package pools

import (
	"runtime"
	"sync"
	"testing"

	"github.com/dgraph-io/ristretto"
"github.com/efficientgo/core/testutil"
bktpool "github.com/gobwas/pool"
)
type client struct {
	innerUpload func(fileName string, chunkBuffer []byte)

	pool         sync.Pool
	bucketedPool *bktpool.Pool
	cache        *ristretto.Cache
}

func getSize(fileName string) int {
	switch fileName {
	case "a.txt":
		return 800
	case "b.txt":
		return 4 * 800
	case "a.800MB.txt":
		return 800e6
	case "b.100GB.txt":
		return 100e9
	}
	panic("unknown file")
}

func (c *client) upload(fileName string) {
	chunkSize := getSize(fileName) / 800

	b := make([]byte, chunkSize)

	c.innerUpload(fileName, b)
}

func (c *client) uploadWithPool(fileName string) {
	chunkSize := getSize(fileName) / 800

	b := c.pool.Get().([]byte)

	if cap(b) < chunkSize {
		b = make([]byte, chunkSize)
	}

	b = b[:chunkSize]
	c.innerUpload(fileName, b)

	c.pool.Put(b)
}

func (c *client) uploadWithBucketedPool(fileName string) {
	chunkSize := getSize(fileName) / 800

	bi, _ := c.bucketedPool.Get(chunkSize)
	b := bi.([]byte)

	c.innerUpload(fileName, b)

	c.bucketedPool.Put(b, cap(b))
}

func (c *client) uploadWithCache(fileName string) {
	chunkSize := getSize(fileName) / 800

	value, found := c.cache.Get(chunkSize)
	var b []byte
	if found {
		b = value.([]byte)
	}

	if cap(b) < chunkSize {
		b = make([]byte, chunkSize)
	}

	c.innerUpload(fileName, b[:chunkSize])

	c.cache.Set(chunkSize, b, 10)
}

func benchmarkUploadSimple(b *testing.B, uploadFn func(string)) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	for g := 0; g < 2; g++ {
		go func() {
			for i := 0; i < b.N; i++ {
				if i == b.N/2 {
					runtime.GC()
				}

				if i%2 == 0 {
					uploadFn("a.800MB.txt")
					continue
				}
				uploadFn("b.100GB.txt")

			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func benchmarkUpload(b *testing.B, uploadFn func(string)) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	workCh := make(chan string)
	b.ResetTimer()
	for w := 0; w < 10; w++ {
		go func() {
			for {
				f, ok := <-workCh
				if !ok {
					wg.Done()
					return
				}
				uploadFn(f)
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		if i == b.N/2 {
			runtime.GC()
		}
		if i%2 == 0 {
			workCh <- "a.800MB.txt"
			continue
		}
		workCh <- "b.100GB.txt"
	}
	close(workCh)
	wg.Wait()
}

// BenchmarkUploads recommended run:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkUploads$' -benchtime 100x -count=5 -cpu 4 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkUploads(b *testing.B) {
	BenchmarkUploads_Make(b)
}

func BenchmarkUploads_Make(b *testing.B) {
	b.ReportAllocs()

	cl := &client{
		innerUpload: func(_ string, b []byte) {
			// Simulate some work that depends on buffer length.
			for i := range b {
				b[i] = 'a'
			}
		},
	}
	benchmarkUpload(b, cl.upload)
}

func BenchmarkUploads_Pool(b *testing.B) {
	b.ReportAllocs()

	cl := &client{
		innerUpload: func(_ string, b []byte) {
			// We have to uniform the latency across uploads, because it's unfair
			// when sync.Pool allocates a bit more which causes benchmark to have time for more
			// GC runs. This skews the results and shows sync.Pool actually allocating more from time to time.
			for i := range b {
				b[i] = 'a'
			}
		},
		pool: sync.Pool{
			New: func() any { return []byte(nil) },
		},
	}
	benchmarkUpload(b, cl.uploadWithPool)
}

func BenchmarkUploads_BucketedPool(b *testing.B) {
	b.ReportAllocs()

	cl := &client{
		innerUpload: func(_ string, b []byte) {
			// We have to uniform the latency across uploads, because it's unfair
			// when sync.Pool allocates a bit more which causes benchmark to have time for more
			// GC runs. This skews the results and shows sync.Pool actually allocating more from time to time.
			for i := range b {
				b[i] = 'a'
			}
		},
		bucketedPool: bktpool.New(1e6, 128e6),
	}
	benchmarkUpload(b, cl.uploadWithBucketedPool)
}

func BenchmarkUploads_Cache(b *testing.B) {
	b.ReportAllocs()

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     500e6,
		BufferItems: 10,
	})
	testutil.Ok(b, err)

	cl := &client{
		innerUpload: func(_ string, b []byte) {
			// We have to uniform the latency across uploads, because it's unfair
			// when sync.Pool allocates a bit more which causes benchmark to have time for more
			// GC runs. This skews the results and shows sync.Pool actually allocating more from time to time.
			for i := range b {
				b[i] = 'a'
			}
		},
		cache: cache,
	}
	benchmarkUpload(b, cl.uploadWithCache)
}

func BenchmarkUploads_StaticBufs(b *testing.B) {
	b.ReportAllocs()

	cl := &client{
		innerUpload: func(_ string, b []byte) {
			// We have to uniform the latency across uploads, because it's unfair
			// when sync.Pool allocates a bit more which causes benchmark to have time for more
			// GC runs. This skews the results and shows sync.Pool actually allocating more from time to time.
			for i := range b {
				b[i] = 'a'
			}
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	for g := 0; g < 2; g++ {
		go func() {
			buf := make([]byte, 128e6)
			for i := 0; i < b.N; i++ {
				if i%2 == 0 {
					cl.innerUpload("a.800MB.txt", buf[:1e6])
					continue
				}
				cl.innerUpload("b.100GB.txt", buf[:128e6])
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestUploads(t *testing.T) {
	cl := &client{}
	for _, tcase := range []struct {
		name     string
		uploadFn func(string)
	}{
		{name: "make", uploadFn: cl.upload},
		{name: "sync-pool", uploadFn: cl.uploadWithPool},
		{name: "bucket-pool", uploadFn: cl.uploadWithBucketedPool},
		{name: "cache", uploadFn: cl.uploadWithCache},
		{name: "static-bufs", uploadFn: nil},
	} {
		if ok := t.Run(tcase.name, func(t *testing.T) {
			messages := map[int][][]byte{}
			var mu sync.Mutex

			cl.innerUpload = func(f string, chunkBuffer []byte) {
				mu.Lock()

				for i := range chunkBuffer {
					chunkBuffer[i] = f[0]
				}

				// Copy as those bytes can be modified in place.
				cb := make([]byte, len(chunkBuffer))
				copy(cb, chunkBuffer)
				messages[len(chunkBuffer)] = append(messages[len(chunkBuffer)], cb)
				mu.Unlock()
			}
			cl.pool.New = func() any { return []byte(nil) }
			cl.bucketedPool = bktpool.New(10, 1e3)

			cache, err := ristretto.NewCache(&ristretto.Config{
				NumCounters: 1e6,
				MaxCost:     5e6,
				BufferItems: 10,
			})
			testutil.Ok(t, err)
			cl.cache = cache

			if tcase.uploadFn != nil {
				tcase.uploadFn("a.txt")
				tcase.uploadFn("b.txt")
				tcase.uploadFn("a.txt")
				tcase.uploadFn("b.txt")
			} else {
				buf := make([]byte, 4)
				cl.innerUpload("a.txt", buf[:1])
				cl.innerUpload("b.txt", buf[:4])
				cl.innerUpload("a.txt", buf[:1])
				cl.innerUpload("b.txt", buf[:4])
			}

			testutil.Equals(t, map[int][][]byte{
				1: {
					[]byte{'a'},
					[]byte{'a'},
				},
				4: {
					[]byte{'b', 'b', 'b', 'b'},
					[]byte{'b', 'b', 'b', 'b'},
				},
			}, messages)
		}); !ok {
			return
		}
	}

}
