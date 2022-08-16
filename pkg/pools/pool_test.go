package pools

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

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
	for i := 0; i < b.N; i++ {

	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	sendFn('b', 1e6)
	sendFn('a', 1e3)

	go func() {
		defer wg.Done()

		for k := 0; k < b.N; k++ {
			sendFn('a', 1e3)
		}
	}()

	go func() {
		defer wg.Done()

		for k := 0; k < b.N; k++ {
			sendFn('b', 1e6)
		}
	}()
	wg.Wait()
}

func BenchmarkSend(b *testing.B) {
	b.ReportAllocs()

	cl := &client{}
	b.ResetTimer()
	benchmarkSend(b, cl, cl.send)
}

func BenchmarkSend2(b *testing.B) {
	b.ReportAllocs()

	cl := &client{}
	cl.pool.New = func() any { return []byte(nil) }
	b.ResetTimer()
	benchmarkSend(b, cl, cl.sendWithPool)

	// Tool that counts memory used by certain structure would be nice....
	m := runtime.MemStats{}
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println(m.HeapAlloc)

	runtime.KeepAlive(cl)

	// BenchmarkSend2
	//1199936
	//1202424
	//1202840
	//BenchmarkSend2-12    	    4317	    282298 ns/op	     282 B/op	       2 allocs/op
	//PASS
}

func BenchmarkSend3(b *testing.B) {
	b.ReportAllocs()

	cl := &client{}
	cl.bucketedPool = NewBucketedPool(1e3, 1e6)
	b.ResetTimer()
	benchmarkSend(b, cl, cl.sendWithBucketedPool)

	// Tool that counts memory used by certain structure would be nice....
	m := runtime.MemStats{}
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println(m.HeapAlloc)

	runtime.KeepAlive(cl)

	//BenchmarkSend3
	//1205976
	//1208368
	//1208784
	//BenchmarkSend3-12    	    2533	    424419 ns/op	     400 B/op	       0 allocs/op
	//PASS
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
