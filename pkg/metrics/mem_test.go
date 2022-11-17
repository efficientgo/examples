package metrics

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"runtime/metrics"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func naivePrintMemStats() {
	runtime.GC()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	// Example output: 2022/04/09 13:42:33 {
	// Alloc:472536
	// TotalAlloc:773208
	// Sys:11027464
	// Lookups:0
	// Mallocs:3543
	// Frees:1929
	// HeapAlloc:472536
	// HeapSys:3735552
	// HeapIdle:2170880
	// HeapInuse:1564672
	// HeapReleased:1720320
	// HeapObjects:1614
	// StackInuse:458752
	// StackSys:458752
	// MSpanInuse:55080
	// MSpanSys:65536
	// MCacheInuse:14400
	// MCacheSys:16384
	// BuckHashSys:1445701
	// GCSys:4009176
	// OtherSys:1296363
	// NextGC:4194304
	// LastGC:1649508153943084326
	// PauseTotalNs:15783
	// PauseNs:[15783 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	// PauseEnd:[1649508153943084326 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	// NumGC:1
	// NumForcedGC:1
	// GCCPUFraction:0.015881668262409922
	// EnableGC:true
	// DebugGC:false
	// BySize:[{Size:0 Mallocs:0 Frees:0} {Size:8 Mallocs:83 Frees:37} {Size:16 Mallocs:897 Frees:457} {Size:24 Mallocs:508 Frees:431} {Size:32 Mallocs:307 Frees:155} {Size:48 Mallocs:302 Frees:80} {Size:64 Mallocs:149 Frees:47} {Size:80 Mallocs:62 Frees:51} {Size:96 Mallocs:88 Frees:35} {Size:112 Mallocs:443 Frees:260} {Size:128 Mallocs:27 Frees:14} {Size:144 Mallocs:0 Frees:0} {Size:160 Mallocs:79 Frees:36} {Size:176 Mallocs:25 Frees:3} {Size:192 Mallocs:1 Frees:0} {Size:208 Mallocs:60 Frees:34} {Size:224 Mallocs:1 Frees:0} {Size:240 Mallocs:1 Frees:0} {Size:256 Mallocs:23 Frees:6} {Size:288 Mallocs:14 Frees:5} {Size:320 Mallocs:35 Frees:26} {Size:352 Mallocs:25 Frees:6} {Size:384 Mallocs:2 Frees:1} {Size:416 Mallocs:60 Frees:15} {Size:448 Mallocs:12 Frees:0} {Size:480 Mallocs:3 Frees:0} {Size:512 Mallocs:2 Frees:2} {Size:576 Mallocs:11 Frees:7} {Size:640 Mallocs:29 Frees:15} {Size:704 Mallocs:8 Frees:3} {Size:768 Mallocs:0 Frees:0} {Size:896 Mallocs:15 Frees:11} {Size:1024 Mallocs:13 Frees:2} {Size:1152 Mallocs:6 Frees:3} {Size:1280 Mallocs:16 Frees:7} {Size:1408 Mallocs:5 Frees:3} {Size:1536 Mallocs:7 Frees:2} {Size:1792 Mallocs:15 Frees:6} {Size:2048 Mallocs:1 Frees:0} {Size:2304 Mallocs:6 Frees:0} {Size:2688 Mallocs:8 Frees:3} {Size:3072 Mallocs:0 Frees:0} {Size:3200 Mallocs:1 Frees:0} {Size:3456 Mallocs:0 Frees:0} {Size:4096 Mallocs:8 Frees:4} {Size:4864 Mallocs:3 Frees:3} {Size:5376 Mallocs:2 Frees:0} {Size:6144 Mallocs:7 Frees:4} {Size:6528 Mallocs:0 Frees:0} {Size:6784 Mallocs:0 Frees:0} {Size:6912 Mallocs:0 Frees:0} {Size:8192 Mallocs:2 Frees:0} {Size:9472 Mallocs:2 Frees:0} {Size:9728 Mallocs:0 Frees:0} {Size:10240 Mallocs:12 Frees:0} {Size:10880 Mallocs:1 Frees:1} {Size:12288 Mallocs:0 Frees:0} {Size:13568 Mallocs:1 Frees:0} {Size:14336 Mallocs:0 Frees:0} {Size:16384 Mallocs:0 Frees:0} {Size:18432 Mallocs:0 Frees:0}]}
	log.Printf("%+v\n", mem)
}

var memMetrics = []metrics.Sample{
	// Total bytes allocated.
	{Name: "/gc/heap/allocs:bytes"},
	// Currently used bytes on heap.
	{Name: "/memory/classes/heap/objects:bytes"},
}

// Example of getting CPU usage through Prometheus metrics
// Read more in "Efficient Go"; Example 6-12.
func printMemRuntimeMetric() {
	runtime.GC()
	metrics.Read(memMetrics)

	fmt.Println("Total bytes allocaed:", memMetrics[0].Value.Uint64())
	fmt.Println("Current inuse bytes:", memMetrics[1].Value.Uint64())
}

func TestRTMetrics(t *testing.T) {
	printMemRuntimeMetric()
	naivePrintMemStats()
}

func ExampleMemoryMetrics() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{
				Matcher: regexp.MustCompile("/gc/heap/allocs:bytes"),
			},
			collectors.GoRuntimeMetricsRule{
				Matcher: regexp.MustCompile("/memory/classes/heap/objects:bytes"),
			},
		)))

	go http.ListenAndServe(":8080", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	for i := 0; i < xTimes; i++ {
		err := doOperation()
		// ...
		_ = err
	}

	printPrometheusMetrics(reg)
}
