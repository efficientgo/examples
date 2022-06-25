package metrics

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"time"

	"github.com/bwplotka/tracing-go/tracing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prepare() { fmt.Println("initializing operation!") }

//nolint
func doOperation() error {
	// Do some dummy, randomized heavy work (both in terms of latency, CPU and memory usage).
	alloc := make([]byte, 1e6)
	for i := 0; i < int(rand.Float64()*100); i++ {
		_ = fmt.Sprintf("doing stuff! %+v", alloc)
	}

	runtime.GC() // To have more interesting GC metrics.

	switch rand.Intn(3) {
	case 0:
		return nil
	case 1:
		return errors.New("error first")
	case 2:
		return errors.New("error other")
	}
	return nil
}

func doOperationWithCtx(ctx context.Context) error {
	_, span := tracing.StartSpan(ctx, "first operation")
	// Do some dummy, randomized heavy work (both in terms of latency, CPU and memory usage).
	alloc := make([]byte, 1e6)
	for i := 0; i < int(rand.Float64()*125); i++ {
		_ = fmt.Sprintf("doing stuff! %+v", alloc) // Hope for this to not get cleared by compiler.
	}
	span.End(nil)

	runtime.GC() // To have more interesting GC metrics.

	_ = tracing.DoInSpan(ctx, "sub operation2", func(ctx context.Context, span tracing.Span) error {
		return nil
	})

	_ = tracing.DoInSpan(ctx, "sub operation3", func(ctx context.Context, span tracing.Span) error {
		return nil
	})

	return tracing.DoInSpan(
		ctx,
		"choosing error",
		func(ctx context.Context, span tracing.Span) error {
			switch rand.Intn(3) {
			default:
				return nil
			case 1:
				time.Sleep(300 * time.Millisecond) // For more interesting results.
				return errors.New("error first")
			case 2:
				return errors.New("error other")
			}
		})
}

func tearDown() { fmt.Println("closing operation!") }

func errorType(err error) string {
	if err != nil {
		if err.Error() == "error first" {
			return "error1"
		}
		return "other_error"
	}
	return ""
}

func printPrometheusMetrics(reg prometheus.Gatherer) {
	rec := httptest.NewRecorder()
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{DisableCompression: true, EnableOpenMetrics: true}).ServeHTTP(rec, &http.Request{})
	if rec.Code != 200 {
		panic("unexpected error")
	}

	fmt.Println(rec.Body.String())
}
