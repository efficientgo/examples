package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Example of getting CPU usage through Prometheus metrics
// Read more in "Efficient Go"; Example 6-11.
func ExampleCPUTimeMetric() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	go http.ListenAndServe(":8080", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	for i := 0; i < xTimes; i++ {
		err := doOperation()
		// ...
		_ = err
	}

	printPrometheusMetrics(reg)
}
