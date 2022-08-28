package metrics

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bwplotka/tracing-go/tracing"
	"github.com/bwplotka/tracing-go/tracing/exporters/jaeger"
	"github.com/bwplotka/tracing-go/tracing/exporters/otlp"
	"github.com/efficientgo/e2e"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
"github.com/efficientgo/core/testutil"
"github.com/go-kit/log"
"github.com/go-kit/log/level"
"github.com/prometheus/client_golang/prometheus"
"github.com/prometheus/client_golang/prometheus/promauto"
"github.com/prometheus/client_golang/prometheus/promhttp"
)

const xTimes = 10

func ExampleLatencySimplest() {
	prepare()

	for i := 0; i < xTimes; i++ {
		start := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(start)

		fmt.Printf("%v ns\n", elapsed.Nanoseconds())

		if err != nil { /* Handle error... */
		}
	}

	tearDown()

	// Output:
}

func ExampleLatencySimplestAggr() {
	var count, sum int64

	prepare()

	for i := 0; i < xTimes; i++ {
		start := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(start)

		sum += elapsed.Nanoseconds()
		count++

		if err != nil { /* Handle error... */
		}
	}

	fmt.Printf("%v ns/op\n", sum/count) // 188324467 ns/op

	tearDown()

	// Output:
}

func ExampleLatencyLog() {
	logger := log.With(log.NewLogfmtLogger(os.Stderr), "ts", log.DefaultTimestampUTC)

	prepare()

	for i := 0; i < xTimes; i++ {
		now := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(now)

		// Log line level=info ts=2022-05-02T11:30:47.803680146Z msg="finished operation" result="error first" elapsed=292.639849ms
		level.Info(logger).Log("msg", "finished operation", "result", err, "elapsed", elapsed.String())

		if err != nil { /* Handle error... */
		}
	}

	tearDown()

	// Output:
}

func ExampleLatencyTrace() {
	tracer, cleanFn, err := tracing.NewTracer(otlp.Exporter("<your tracing collector>"))
	if err != nil { /* Handle error... */
	}
	defer cleanFn()

	prepare()

	for i := 0; i < xTimes; i++ {
		ctx, span := tracer.StartSpan("doOperation")
		err := doOperationWithCtx(ctx) // Operation we want to measure and potentially optimize...
		span.End(err)

		if err != nil { /* Handle error... */
		}
	}

	tearDown()

	// Output:
}

func ExampleLatencyMetric() {
	reg := prometheus.NewRegistry()
	latencySeconds := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "operation_duration_seconds",
		Help:    "Tracks the latency of operations in seconds.",
		Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
	}, []string{"error_type"})

	prepare()

	go func() {
		for i := 0; i < xTimes; i++ {
			now := time.Now()
			err := doOperation() // Operation we want to measure and potentially optimize...
			elapsed := time.Since(now)

			// Prometheus metric.
			latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())

			if err != nil { /* Handle error... */
			}

			time.Sleep(1 * time.Second)
		}
	}()

	if err := http.ListenAndServe(":8080", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})); err != nil {
		stdlog.Fatal(err)
	}

	tearDown()

	printPrometheusMetrics(reg)

	// Output:
}

func TestLatencyE2e(t *testing.T) {
	e, err := e2e.NewDockerEnvironment("e2e_latency")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	reg := prometheus.NewRegistry()
	mon, err := e2emonitoring.Start(e,
		e2emonitoring.WithScrapeInterval(1*time.Second),
		e2emonitoring.WithCustomRegistry(reg),
	)
	testutil.Ok(t, err)

	// Setup in-memory Jaeger to check if backend can understand our client traces.
	j := e.Runnable("tracing").
		WithPorts(
			map[string]int{
				"http.front":    16686,
				"jaeger.thrift": 16000,
			}).
		Init(e2e.StartOptions{
			Image:   "jaegertracing/all-in-one:1.33",
			Command: e2e.NewCommand("--collector.http-server.host-port=:16000"),
		})

	testutil.Ok(t, e2e.StartAndWaitReady(j))

	tracer, cleanFn, err := tracing.NewTracer(
		jaeger.Exporter("http://"+j.Endpoint("jaeger.thrift")+"/api/traces"),
		tracing.WithServiceName("example"),
	)
	testutil.Ok(t, err)
	t.Cleanup(func() { cleanFn() })

	logger := log.With(log.NewLogfmtLogger(os.Stderr), "ts", log.DefaultTimestampUTC)
	latencySeconds := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "operation_duration_seconds",
		Help:    "Tracks the latency of operations in seconds.",
		Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
	}, []string{"error_type"})

	prepare()

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for ctx.Err() == nil {
			time.Sleep(50 * time.Millisecond)

			// Instrumentation COMBO!
			sctx, span := tracer.StartSpan("doOperation", tracing.WithTracerStartSpanContext(ctx))
			start := time.Now()
			err := doOperationWithCtx(sctx) // Operation we want to measure and potentially optimize...
			elapsed := time.Since(start)
			span.End(err)

			level.Info(logger).Log("msg", "finished operation", "result", err, "elapsed", elapsed.String())

			if span.Context().IsSampled() {
				latencySeconds.WithLabelValues(errorType(err)).(prometheus.ExemplarObserver).ObserveWithExemplar(
					elapsed.Seconds(), map[string]string{"trace-id": span.Context().TraceID()})
			} else {
				latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())
			}

			// Handle error...
			if err != nil {
			}
		}
		wg.Done()
	}()

	// TODO(bwplotka): Make it non-interactive and expect certain Jaeger output.
	testutil.Ok(t, e2einteractive.OpenInBrowser("http://"+j.Endpoint("http.front")))
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser("/graph?g0.expr=rate(operation_duration_seconds_sum%5B1m%5D)%20%2F%20rate(operation_duration_seconds_count%5B1m%5D)&g0.tab=0&g0.stacked=0&g0.range_input=10m&g0.end_input=2022-04-09%2020%3A20%3A40&g0.moment_input=2022-04-09%2020%3A20%3A40&g1.expr=histogram_quantile(0.9%2C%20sum%20by(error_type%2C%20le)%20(rate(operation_duration_seconds_bucket%5B1m%5D)))&g1.tab=0&g1.stacked=0&g1.range_input=10m&g1.end_input=2022-04-09%2020%3A20%3A40&g1.moment_input=2022-04-09%2020%3A20%3A40&g2.expr=operation_duration_seconds_bucket&g2.tab=1&g2.stacked=0&g2.range_input=1h&g3.expr=delta(operation_duration_seconds_count%5B1m%5D)&g3.tab=0&g3.stacked=0&g3.range_input=15m"))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())

	cancel()
	wg.Wait()

	tearDown()
}

func BenchmarkExampleLatency(b *testing.B) {
	prepare()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doOperation()
	}
}

var latTest time.Duration

func BenchmarkLatencyItself(b *testing.B) {
	for i := 0; i < b.N; i++ {
		start := time.Now()
		latTest = time.Since(start)
	}
}
