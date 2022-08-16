package pools

import (
	"sync"
	"testing"

	"github.com/dgraph-io/ristretto"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

type client struct {
	innerUpload func(fileName string, chunkBuffer []byte)

	pool         sync.Pool
	bucketedPool *BucketedPool
	cache        *ristretto.Cache
}

func (c *client) upload(fileName string, chunkSize int) {
	b := make([]byte, chunkSize)

	c.innerUpload(fileName, b)
}

func (c *client) uploadWithPool(fileName string, chunkSize int) {
	b := c.pool.Get().([]byte)

	if cap(b) < chunkSize {
		b = make([]byte, chunkSize)
	}

	b = b[:chunkSize]
	c.innerUpload(fileName, b)

	c.pool.Put(b)
}

func (c *client) uploadWithBucketedPool(fileName string, chunkSize int) {
	b := c.bucketedPool.Get(chunkSize)

	c.innerUpload(fileName, *b)

	c.bucketedPool.Put(b)
}

func (c *client) uploadWithCache(fileName string, chunkSize int) {
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

func benchmarkUpload(b *testing.B, uploadFn func(string, int)) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			uploadFn("a.60KB.txt", 10)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			uploadFn("b.100MB.txt", 1e6)
		}
		wg.Done()
	}()

	wg.Wait()
}

func benchmarkUpload2(b *testing.B, cl *client) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	go func() {
		buf := make([]byte, 10)
		for i := 0; i < b.N; i++ {
			cl.innerUpload("a.60KB.txt", buf)
		}
		wg.Done()
	}()
	go func() {
		buf := make([]byte, 1e6)
		for i := 0; i < b.N; i++ {
			cl.innerUpload("b.100MB.txt", buf)
		}
		wg.Done()
	}()

	wg.Wait()
}

// BenchmarkUploads recommended run:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkUploads' -benchtime 10000x -cpu 4 -benchmem -memprofile=${ver}.mem.pprof | tee ${ver}.txt
func BenchmarkUploads(b *testing.B) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     5e6,
		BufferItems: 10,
	})
	testutil.Ok(b, err)

	cl := &client{
		innerUpload: func(string, []byte) {},
		pool: sync.Pool{
			New: func() any { return []byte(nil) },
		},
		bucketedPool: NewBucketedPool(1e3, 1e6),
		cache:        cache,
	}

	for _, tcase := range []struct {
		name     string
		uploadFn func(string, int)
	}{
		{name: "make", uploadFn: cl.upload},
		{name: "sync-pool", uploadFn: cl.uploadWithPool},
		{name: "bucket-pool", uploadFn: cl.uploadWithBucketedPool},
		{name: "cache", uploadFn: cl.uploadWithCache},
		{name: "static-bufs", uploadFn: nil},
	} {
		b.Run(tcase.name, func(b *testing.B) {
			b.ReportAllocs()

			if tcase.uploadFn != nil {
				benchmarkUpload(b, tcase.uploadFn)
				return
			}

			benchmarkUpload2(b, cl)
		})
	}
}

func TestUploads(t *testing.T) {
	cl := &client{}
	for _, tcase := range []struct {
		name     string
		uploadFn func(string, int)
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
			cl.bucketedPool = NewBucketedPool(10, 1e3)

			cache, err := ristretto.NewCache(&ristretto.Config{
				NumCounters: 1e6,
				MaxCost:     5e6,
				BufferItems: 10,
			})
			testutil.Ok(t, err)
			cl.cache = cache

			if tcase.uploadFn != nil {
				tcase.uploadFn("a.txt", 1)
				tcase.uploadFn("b.txt", 4)
				tcase.uploadFn("a.txt", 1)
				tcase.uploadFn("b.txt", 4)
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
