package pools

import (
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

type client struct {
	forwardFn func([]byte)

	pool         sync.Pool
	bucketedPool *BucketedPool
}

func (c *client) send(char byte, lenToSend int) {
	b := make([]byte, lenToSend)
	for i := range b {
		b[i] = char
	}
	c.forwardFn(b)
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
	c.forwardFn(b)

	c.pool.Put(b)
}

func (c *client) sendWithBucketedPool(char byte, lenToSend int) {
	b := c.bucketedPool.Get(lenToSend)

	for i := range *b {
		(*b)[i] = char
	}
	c.forwardFn(*b)

	c.bucketedPool.Put(b)
}

func benchmarkSend(b *testing.B, cl *client, sendFn func(byte, int)) {
	wg := sync.WaitGroup{}
	wg.Add(2 * b.N)
	cl.forwardFn = func([]byte) {
		time.Sleep(50 * time.Millisecond)

		wg.Done()
	}

	for i := 0; i < b.N; i++ {
		go sendFn('a', 1e3)
		go sendFn('b', 1e6)

		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
}

// BenchmarkSends recommended run:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkSends' -benchtime 4000x -cpu 4 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkSends(b *testing.B) {
	cl := &client{}
	cl.pool.New = func() any { return []byte(nil) }
	cl.bucketedPool = NewBucketedPool(1e3, 1e6)

	for _, tcase := range []struct {
		name   string
		sendFn func(byte, int)
	}{
		{name: "make", sendFn: cl.send},
		{name: "sync-pool", sendFn: cl.sendWithPool},
		{name: "bucket-pool", sendFn: cl.sendWithBucketedPool},
	} {
		b.Run(tcase.name, func(b *testing.B) {
			b.ReportAllocs()

			b.ResetTimer()
			benchmarkSend(b, cl, tcase.sendFn)
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
	} {
		if ok := t.Run("", func(t *testing.T) {
			messages := map[int][][]byte{}
			var mu sync.Mutex

			cl.forwardFn = func(b []byte) {
				mu.Lock()

				// Copy as those bytes can be modified in place.
				cb := make([]byte, len(b))
				copy(cb, b)
				messages[len(b)] = append(messages[len(b)], cb)
				mu.Unlock()
			}
			cl.pool.New = func() any { return []byte(nil) }
			cl.bucketedPool = NewBucketedPool(10, 1e3)

			tcase.sendFn('a', 1)
			tcase.sendFn('b', 4)
			tcase.sendFn('a', 1)
			tcase.sendFn('b', 4)

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
