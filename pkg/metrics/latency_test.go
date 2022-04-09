package metrics

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prepare() { fmt.Println("initializing operation!") }

//nolint
func doOperation() error {
	// Do some dummy, heavy work (both in terms of latency, CPU and memory usage).
	alloc := make([]byte, 1e6)
	for i := 0; i < int(rand.Float64()*100); i++ {
		_ = fmt.Sprintf("doing stuff! %+v", alloc)
	}

	switch rand.Intn(3) {
	case 0:
		return nil
	case 1:
		return errors.New("err1")
	case 2:
		return errors.New("other error")
	}
	return nil
}

func tearDown() { fmt.Println("closing operation!") }

func errorType(err error) string {
	if err != nil {
		if err.Error() == "err1" {
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

func ExampleLatency() {
	logger := log.NewLogfmtLogger(os.Stderr)
	reg := prometheus.NewRegistry()
	latencySeconds := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "operation_duration_seconds",
		Help:    "Tracks the latency of operations in seconds.",
		Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
	}, []string{"error_type"})

	prepare()

	now := time.Now()
	err := doOperation() // Operation we want to measure and potentially optimize...
	elapsed := time.Since(now)

	// Log line level=info msg="finished operation" result=err1 elapsed=194.715965ms
	level.Info(logger).Log("msg", "finished operation", "result", err, "elapsed", elapsed.String())
	// Prometheus metric.
	latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())

	// Handle error...
	if err != nil {
	}

	tearDown()

	printPrometheusMetrics(reg)

	// Output:
	// initializing operation!
	// closing operation!
	// # HELP operation_duration_seconds Tracks the latency of operations in seconds.
	// # TYPE operation_duration_seconds histogram
	// operation_duration_seconds_bucket{error_type="",le="0.001"} 0
	// operation_duration_seconds_bucket{error_type="",le="0.01"} 0
	// operation_duration_seconds_bucket{error_type="",le="0.1"} 0
	// operation_duration_seconds_bucket{error_type="",le="1"} 0
	// operation_duration_seconds_bucket{error_type="",le="10"} 1
	// operation_duration_seconds_bucket{error_type="",le="100"} 1
	// operation_duration_seconds_bucket{error_type="",le="+Inf"} 1
	// operation_duration_seconds_sum{error_type=""} 1.416099937
	// operation_duration_seconds_count{error_type=""} 1
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

	logger := log.NewLogfmtLogger(os.Stderr)
	latencySeconds := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "operation_duration_seconds",
		Help:    "Tracks the latency of operations in seconds.",
		Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
	}, []string{"error_type"})

	prepare()

	for i := 0; i < 100; i++ {
		now := time.Now()
		err := doOperation() // Operation we want to measure and potentially optimize...
		elapsed := time.Since(now)

		// Log line, example level=info msg="finished operation" result=err1 elapsed=381.832888ms
		level.Info(logger).Log("msg", "finished operation", "result", err, "elapsed", elapsed.String())
		// Prometheus metric.
		latencySeconds.WithLabelValues(errorType(err)).Observe(elapsed.Seconds())

		// Handle error...
		if err != nil {
		}
	}
	tearDown()

	testutil.Ok(t, mon.OpenUserInterfaceInBrowser("/graph?g0.expr=rate(operation_duration_seconds_sum%5B1m%5D)%20%2F%20rate(operation_duration_seconds_count%5B1m%5D)&g0.tab=0&g0.stacked=0&g0.range_input=2m&g0.end_input=2022-04-09%2020%3A20%3A40&g0.moment_input=2022-04-09%2020%3A20%3A40&g1.expr=histogram_quantile(0.9%2C%20sum%20by(error_type%2C%20le)%20(rate(operation_duration_seconds_bucket%5B1m%5D)))&g1.tab=0&g1.stacked=0&g1.range_input=1m&g1.end_input=2022-04-09%2020%3A20%3A40&g1.moment_input=2022-04-09%2020%3A20%3A40&g2.expr=operation_duration_seconds_bucket&g2.tab=1&g2.stacked=0&g2.range_input=1h&g3.expr=delta(operation_duration_seconds_count%5B1m%5D)&g3.tab=0&g3.stacked=0&g3.range_input=15m"))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}
