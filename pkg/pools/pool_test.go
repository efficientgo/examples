package pools

import (
	"sync"
	"testing"

	"github.com/dgraph-io/ristretto"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

type client struct {
	forward func([]byte)

	pool         sync.Pool
	bucketedPool *BucketedPool
	cache        *ristretto.Cache
}

func (c *client) send(char byte, lenToSend int) {
	b := make([]byte, lenToSend)
	for i := range b {
		b[i] = char
	}
	c.forward(b)
}

func (c *client) sendWithPool(char byte, lenToSend int) {
	b := c.pool.Get().([]byte)

	if cap(b) < lenToSend {
		b = make([]byte, lenToSend)
	}

	b = b[:lenToSend]
	for i := range b {
		b[i] = char
	}
	c.forward(b)

	c.pool.Put(b)
}

func (c *client) sendWithBucketedPool(char byte, lenToSend int) {
	b := c.bucketedPool.Get(lenToSend)

	for i := range *b {
		(*b)[i] = char
	}
	c.forward(*b)

	c.bucketedPool.Put(b)
}

func (c *client) sendWith(b []byte, char byte, lenToSend int) {
	if cap(b) < lenToSend {
		b = make([]byte, lenToSend)
	}

	b = b[:lenToSend]
	for i := range b {
		b[i] = char
	}
	c.forward(b)
}

func (c *client) sendWithCache(char byte, lenToSend int) {
	value, found := c.cache.Get(lenToSend)
	var b []byte
	if found {
		b = value.([]byte)
	}

	if cap(b) < lenToSend {
		b = make([]byte, lenToSend)
	}

	b = b[:lenToSend]
	for i := range b {
		b[i] = char
	}
	c.forward(b)

	c.cache.Set(lenToSend, b, 10)
}

func benchmarkSend(b *testing.B, sendFn func(byte, int)) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			sendFn('a', 10)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			sendFn('b', 1e6)
		}
		wg.Done()
	}()

	wg.Wait()
}

func benchmarkSend2(b *testing.B, cl *client) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b.ResetTimer()
	go func() {
		buf := make([]byte, 10)
		for i := 0; i < b.N; i++ {
			cl.sendWith(buf, 'a', 10)
		}
		wg.Done()
	}()
	go func() {
		buf := make([]byte, 1e6)
		for i := 0; i < b.N; i++ {
			cl.sendWith(buf, 'b', 1e6)
		}
		wg.Done()
	}()

	wg.Wait()
}

// BenchmarkSends recommended run:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkSends' -benchtime 4000x -cpu 4 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkSends(b *testing.B) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     5e6,
		BufferItems: 10,
	})
	testutil.Ok(b, err)

	cl := &client{
		forward: func(b []byte) {},
		pool: sync.Pool{
			New: func() any { return []byte(nil) },
		},
		bucketedPool: NewBucketedPool(1e3, 1e6),
		cache:        cache,
	}

	for _, tcase := range []struct {
		name   string
		sendFn func(byte, int)
	}{
		{name: "make", sendFn: cl.send},
		{name: "sync-pool", sendFn: cl.sendWithPool},
		{name: "bucket-pool", sendFn: cl.sendWithBucketedPool},
		{name: "cache", sendFn: cl.sendWithCache},
		{name: "static-bufs", sendFn: nil},
	} {
		b.Run(tcase.name, func(b *testing.B) {
			b.ReportAllocs()

			if tcase.sendFn != nil {
				benchmarkSend(b, tcase.sendFn)
				return
			}

			benchmarkSend2(b, cl)
		})
	}
}

func TestSends(t *testing.T) {
	cl := &client{}
	for _, tcase := range []struct {
		sendFn func(byte, int)
	}{
		{sendFn: cl.send},
		{sendFn: cl.sendWithPool},
		{sendFn: cl.sendWithBucketedPool},
		{sendFn: cl.sendWithCache},
		{sendFn: nil},
	} {
		if ok := t.Run("", func(t *testing.T) {
			messages := map[int][][]byte{}
			var mu sync.Mutex

			cl.forward = func(b []byte) {
				mu.Lock()

				// Copy as those bytes can be modified in place.
				cb := make([]byte, len(b))
				copy(cb, b)
				messages[len(b)] = append(messages[len(b)], cb)
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

			if tcase.sendFn != nil {
				tcase.sendFn('a', 1)
				tcase.sendFn('b', 4)
				tcase.sendFn('a', 1)
				tcase.sendFn('b', 4)
			} else {
				buf := make([]byte, 4)
				cl.sendWith(buf, 'a', 1)
				cl.sendWith(buf, 'b', 4)
				cl.sendWith(buf, 'a', 1)
				cl.sendWith(buf, 'b', 4)
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
