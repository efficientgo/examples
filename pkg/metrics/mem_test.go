package metrics

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"testing"

	"github.com/efficientgo/tools/performance/pkg/profiles"
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

func printMemRuntimeMetric() {
	runtime.GC()
	metrics.Read(memMetrics)

	log.Printf("Total bytes allocaed: %v\n", memMetrics[0].Value.Uint64())
	log.Printf("Current inuse bytes: %v\n", memMetrics[1].Value.Uint64())
}

func TestRTMetrics(t *testing.T) {
	printMemRuntimeMetric()
	naivePrintMemStats()
}

func ExampleMemUsage() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())

	prepare()

	naivePrintMemStats()
	var err error
	pprof.Do(context.Background(), pprof.Labels("id", "my_operation"), func(ctx context.Context) { // https://github.com/polarsignals/pprof-labels-example/blob/60accf8b4fbebcd5f96b3743663af5745ef74596/printprofile.go
		err = doOperation() // Operation we want to measure and potentially optimize...
	})
	naivePrintMemStats()

	// .https://share.polarsignals.com/3bcc303/
	profiles.Heap("./pprof")

	// Handle error...
	if err != nil {
	}

	tearDown()

	printPrometheusMetrics(reg)

	// Output:
	// initializing operation!
	// doing stuff!
	// closing operation!
	// # HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
	// # TYPE go_gc_cycles_automatic_gc_cycles_total counter
	// go_gc_cycles_automatic_gc_cycles_total 0
	// # HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
	// # TYPE go_gc_cycles_forced_gc_cycles_total counter
	// go_gc_cycles_forced_gc_cycles_total 0
	// # HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
	// # TYPE go_gc_cycles_total_gc_cycles_total counter
	// go_gc_cycles_total_gc_cycles_total 0
	// # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
	// # TYPE go_gc_duration_seconds summary
	// go_gc_duration_seconds{quantile="0"} 0
	// go_gc_duration_seconds{quantile="0.25"} 0
	// go_gc_duration_seconds{quantile="0.5"} 0
	// go_gc_duration_seconds{quantile="0.75"} 0
	// go_gc_duration_seconds{quantile="1"} 0
	// go_gc_duration_seconds_sum 0
	// go_gc_duration_seconds_count 0
	// # HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
	// # TYPE go_gc_heap_allocs_by_size_bytes_total histogram
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 3072
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 7678
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 9936
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 11516
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 12364
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 12637
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 12772
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 12840
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 12884
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 12916
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 12918
	// go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 12921
	// go_gc_heap_allocs_by_size_bytes_total_sum 1.68584e+06
	// go_gc_heap_allocs_by_size_bytes_total_count 12921
	// # HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
	// # TYPE go_gc_heap_allocs_bytes_total counter
	// go_gc_heap_allocs_bytes_total 1.68584e+06
	// # HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
	// # TYPE go_gc_heap_allocs_objects_total counter
	// go_gc_heap_allocs_objects_total 12921
	// # HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
	// # TYPE go_gc_heap_frees_by_size_bytes_total histogram
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 0
	// go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 0
	// go_gc_heap_frees_by_size_bytes_total_sum 0
	// go_gc_heap_frees_by_size_bytes_total_count 0
	// # HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
	// # TYPE go_gc_heap_frees_bytes_total counter
	// go_gc_heap_frees_bytes_total 0
	// # HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
	// # TYPE go_gc_heap_frees_objects_total counter
	// go_gc_heap_frees_objects_total 0
	// # HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
	// # TYPE go_gc_heap_goal_bytes gauge
	// go_gc_heap_goal_bytes 4.473924e+06
	// # HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
	// # TYPE go_gc_heap_objects_objects gauge
	// go_gc_heap_objects_objects 12921
	// # HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
	// # TYPE go_gc_heap_tiny_allocs_objects_total counter
	// go_gc_heap_tiny_allocs_objects_total 0
	// # HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
	// # TYPE go_gc_pauses_seconds_total histogram
	// go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
	// go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
	// go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
	// go_gc_pauses_seconds_total_bucket{le="1.2799999999999998e-07"} 0
	// go_gc_pauses_seconds_total_bucket{le="1.2799999999999998e-06"} 0
	// go_gc_pauses_seconds_total_bucket{le="1.6383999999999998e-05"} 0
	// go_gc_pauses_seconds_total_bucket{le="0.00016383999999999998"} 0
	// go_gc_pauses_seconds_total_bucket{le="0.0020971519999999997"} 0
	// go_gc_pauses_seconds_total_bucket{le="0.020971519999999997"} 0
	// go_gc_pauses_seconds_total_bucket{le="0.26843545599999996"} 0
	// go_gc_pauses_seconds_total_bucket{le="+Inf"} 0
	// go_gc_pauses_seconds_total_sum NaN
	// go_gc_pauses_seconds_total_count 0
	// # HELP go_goroutines Number of goroutines that currently exist.
	// # TYPE go_goroutines gauge
	// go_goroutines 4
	// # HELP go_info Information about the Go environment.
	// # TYPE go_info gauge
	// go_info{version="go1.17"} 1
	// # HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
	// # TYPE go_memory_classes_heap_free_bytes gauge
	// go_memory_classes_heap_free_bytes 0
	// # HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
	// # TYPE go_memory_classes_heap_objects_bytes gauge
	// go_memory_classes_heap_objects_bytes 1.68584e+06
	// # HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
	// # TYPE go_memory_classes_heap_released_bytes gauge
	// go_memory_classes_heap_released_bytes 2.162688e+06
	// # HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
	// # TYPE go_memory_classes_heap_stacks_bytes gauge
	// go_memory_classes_heap_stacks_bytes 327680
	// # HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
	// # TYPE go_memory_classes_heap_unused_bytes gauge
	// go_memory_classes_heap_unused_bytes 18096
	// # HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
	// # TYPE go_memory_classes_metadata_mcache_free_bytes gauge
	// go_memory_classes_metadata_mcache_free_bytes 1984
	// # HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
	// # TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
	// go_memory_classes_metadata_mcache_inuse_bytes 14400
	// # HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
	// # TYPE go_memory_classes_metadata_mspan_free_bytes gauge
	// go_memory_classes_metadata_mspan_free_bytes 11480
	// # HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
	// # TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
	// go_memory_classes_metadata_mspan_inuse_bytes 37672
	// # HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
	// # TYPE go_memory_classes_metadata_other_bytes gauge
	// go_memory_classes_metadata_other_bytes 3.548888e+06
	// # HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
	// # TYPE go_memory_classes_os_stacks_bytes gauge
	// go_memory_classes_os_stacks_bytes 0
	// # HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
	// # TYPE go_memory_classes_other_bytes gauge
	// go_memory_classes_other_bytes 922115
	// # HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
	// # TYPE go_memory_classes_profiling_buckets_bytes gauge
	// go_memory_classes_profiling_buckets_bytes 1.444909e+06
	// # HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
	// # TYPE go_memory_classes_total_bytes gauge
	// go_memory_classes_total_bytes 1.0175752e+07
	// # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
	// # TYPE go_memstats_alloc_bytes gauge
	// go_memstats_alloc_bytes 1.68584e+06
	// # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
	// # TYPE go_memstats_alloc_bytes_total counter
	// go_memstats_alloc_bytes_total 1.68584e+06
	// # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
	// # TYPE go_memstats_buck_hash_sys_bytes gauge
	// go_memstats_buck_hash_sys_bytes 1.444909e+06
	// # HELP go_memstats_frees_total Total number of frees.
	// # TYPE go_memstats_frees_total counter
	// go_memstats_frees_total 0
	// # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
	// # TYPE go_memstats_gc_cpu_fraction gauge
	// go_memstats_gc_cpu_fraction 0
	// # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
	// # TYPE go_memstats_gc_sys_bytes gauge
	// go_memstats_gc_sys_bytes 3.548888e+06
	// # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
	// # TYPE go_memstats_heap_alloc_bytes gauge
	// go_memstats_heap_alloc_bytes 1.68584e+06
	// # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
	// # TYPE go_memstats_heap_idle_bytes gauge
	// go_memstats_heap_idle_bytes 2.162688e+06
	// # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
	// # TYPE go_memstats_heap_inuse_bytes gauge
	// go_memstats_heap_inuse_bytes 1.703936e+06
	// # HELP go_memstats_heap_objects Number of allocated objects.
	// # TYPE go_memstats_heap_objects gauge
	// go_memstats_heap_objects 12921
	// # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
	// # TYPE go_memstats_heap_released_bytes gauge
	// go_memstats_heap_released_bytes 2.162688e+06
	// # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
	// # TYPE go_memstats_heap_sys_bytes gauge
	// go_memstats_heap_sys_bytes 3.866624e+06
	// # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
	// # TYPE go_memstats_last_gc_time_seconds gauge
	// go_memstats_last_gc_time_seconds 0
	// # HELP go_memstats_lookups_total Total number of pointer lookups.
	// # TYPE go_memstats_lookups_total counter
	// go_memstats_lookups_total 0
	// # HELP go_memstats_mallocs_total Total number of mallocs.
	// # TYPE go_memstats_mallocs_total counter
	// go_memstats_mallocs_total 12921
	// # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
	// # TYPE go_memstats_mcache_inuse_bytes gauge
	// go_memstats_mcache_inuse_bytes 14400
	// # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
	// # TYPE go_memstats_mcache_sys_bytes gauge
	// go_memstats_mcache_sys_bytes 16384
	// # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
	// # TYPE go_memstats_mspan_inuse_bytes gauge
	// go_memstats_mspan_inuse_bytes 37672
	// # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
	// # TYPE go_memstats_mspan_sys_bytes gauge
	// go_memstats_mspan_sys_bytes 49152
	// # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
	// # TYPE go_memstats_next_gc_bytes gauge
	// go_memstats_next_gc_bytes 4.473924e+06
	// # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
	// # TYPE go_memstats_other_sys_bytes gauge
	// go_memstats_other_sys_bytes 922115
	// # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
	// # TYPE go_memstats_stack_inuse_bytes gauge
	// go_memstats_stack_inuse_bytes 327680
	// # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
	// # TYPE go_memstats_stack_sys_bytes gauge
	// go_memstats_stack_sys_bytes 327680
	// # HELP go_memstats_sys_bytes Number of bytes obtained from system.
	// # TYPE go_memstats_sys_bytes gauge
	// go_memstats_sys_bytes 1.0175752e+07
	// # HELP go_sched_goroutines_goroutines Count of live goroutines.
	// # TYPE go_sched_goroutines_goroutines gauge
	// go_sched_goroutines_goroutines 4
	// # HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
	// # TYPE go_sched_latencies_seconds histogram
	// go_sched_latencies_seconds_bucket{le="-5e-324"} 0
	// go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 3
	// go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 3
	// go_sched_latencies_seconds_bucket{le="1.2799999999999998e-07"} 3
	// go_sched_latencies_seconds_bucket{le="1.2799999999999998e-06"} 6
	// go_sched_latencies_seconds_bucket{le="1.6383999999999998e-05"} 9
	// go_sched_latencies_seconds_bucket{le="0.00016383999999999998"} 12
	// go_sched_latencies_seconds_bucket{le="0.0020971519999999997"} 12
	// go_sched_latencies_seconds_bucket{le="0.020971519999999997"} 12
	// go_sched_latencies_seconds_bucket{le="0.26843545599999996"} 12
	// go_sched_latencies_seconds_bucket{le="+Inf"} 12
	// go_sched_latencies_seconds_sum NaN
	// go_sched_latencies_seconds_count 12
	// # HELP go_threads Number of OS threads created.
	// # TYPE go_threads gauge
	// go_threads 7
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

	// Output:

}
