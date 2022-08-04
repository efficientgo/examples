package profile

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func tenGoroutinesWaitingForChannel(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-ctx.Done()
			wg.Done()
		}()
	}
}

func fiveGoroutinesLocked(l *sync.Mutex, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			l.Lock()
			l.Unlock()
			wg.Done()
		}()
	}
}

func TestGoroutines(t *testing.T) {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	tenGoroutinesWaitingForChannel(ctx, &wg)

	var l sync.Mutex
	l.Lock()
	fiveGoroutinesLocked(&l, &wg)

	defer func() {
		l.Unlock()
		cancel()
		wg.Wait()
	}()

	f, err := os.Create("goroutine.pprof")
	testutil.Ok(t, err)
	testutil.Ok(t, pprof.Lookup("goroutine").WriteTo(f, 0))
}

func TestReceiveFromEmptyCh(t *testing.T) {
	ch := make(chan error, 10)

	fmt.Println(<-ch)
}