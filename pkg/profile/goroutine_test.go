package profile

import (
	"context"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func funcA(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			funcB(ctx, wg)
			<-ctx.Done()
			wg.Done()
		}()
	}
	funcB(ctx, wg)

	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		<-ctx.Done()
		wg.Done()
	}()
}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			<-ctx.Done()
			wg.Done()
		}()
	}
}

func TestGoroutines(t *testing.T) {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	funcA(ctx, &wg)

	f, err := os.Create("goroutine.pprof")
	testutil.Ok(t, err)
	testutil.Ok(t, pprof.Lookup("goroutine").WriteTo(f, 0))

	cancel()
	wg.Wait()
}
