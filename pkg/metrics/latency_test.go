package metrics

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bwplotka/tracing-go/tracing"
	"github.com/bwplotka/tracing-go/tracing/exporters/jaeger"
	"github.com/bwplotka/tracing-go/tracing/exporters/otlp"
	"github.com/efficientgo/e2e"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const xTimes = 10

func ExampleLatencySimplest() {
	prepare()

	for i := 0; i < xTimes; i++ {
		start := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(start)

		fmt.Printf("%vns\n", elapsed.Nanoseconds())

		if err != nil { /* Handle error... */
		}
	}

	tearDown()

	// Output:
}

func ExampleLatencyLog() {
	logger := log.NewLogfmtLogger(os.Stderr)

	prepare()

	for i := 0; i < xTimes; i++ {
		now := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(now)

		// Log line level=info msg="finished operation" result=err1 elapsed=194.715965ms
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
		_, span := tracer.StartSpan("doOperation")
		err := doOperation() // Operation we want to measure and potentially optimize...
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

	for i := 0; i < xTimes; i++ {
		now := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(now)

		// Prometheus metric.
		latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())

		if err != nil { /* Handle error... */
		}
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

	tracer, cleanFn, err := tracing.NewTracer(jaeger.Exporter("http://" + j.Endpoint("jaeger.thrift") + "/api/traces"))
	testutil.Ok(t, err)
	t.Cleanup(func() { cleanFn() })

	logger := log.NewLogfmtLogger(os.Stderr)
	latencySeconds := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "operation_duration_seconds",
		Help:    "Tracks the latency of operations in seconds.",
		Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
	}, []string{"error_type"})

	prepare()

	for i := 0; i < 100; i++ {
		// Instrumentation COMBO!
		_, span := tracer.StartSpan("doOperation")
		start := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(start)
		span.End(err)

		// Log line, example level=info msg="finished operation" result=err1 elapsed=381.832888ms
		level.Info(logger).Log("msg", "finished operation", "result", err, "elapsed", elapsed.String())
		// Prometheus metric.
		latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())

		// Handle error...
		if err != nil {
		}
	}
	tearDown()

	// TODO(bwplotka): Make it non-interactive and expect certain Jaeger output.
	testutil.Ok(t, e2einteractive.OpenInBrowser("http://"+j.Endpoint("http.front")))
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser("/graph?g0.expr=rate(operation_duration_seconds_sum%5B1m%5D)%20%2F%20rate(operation_duration_seconds_count%5B1m%5D)&g0.tab=0&g0.stacked=0&g0.range_input=10m&g0.end_input=2022-04-09%2020%3A20%3A40&g0.moment_input=2022-04-09%2020%3A20%3A40&g1.expr=histogram_quantile(0.9%2C%20sum%20by(error_type%2C%20le)%20(rate(operation_duration_seconds_bucket%5B1m%5D)))&g1.tab=0&g1.stacked=0&g1.range_input=10m&g1.end_input=2022-04-09%2020%3A20%3A40&g1.moment_input=2022-04-09%2020%3A20%3A40&g2.expr=operation_duration_seconds_bucket&g2.tab=1&g2.stacked=0&g2.range_input=1h&g3.expr=delta(operation_duration_seconds_count%5B1m%5D)&g3.tab=0&g3.stacked=0&g3.range_input=15m"))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}

var errTest error

func BenchmarkExampleLatency(b *testing.B) {
	prepare()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errTest = doOperation()
	}
}

var latTest time.Duration

func BenchmarkLatencyItself(b *testing.B) {
	for i := 0; i < b.N; i++ {
		start := time.Now()
		latTest = time.Since(start)
	}
}
