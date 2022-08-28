package dup

import (
	"sort"
	"testing"

"github.com/efficientgo/core/testutil"
)

func TestDeduplicate(t *testing.T) {
	for _, funcCase := range []struct {
		name      string
		dedupFunc func([]int) []int
	}{
		{name: "fast", dedupFunc: DeduplicateFaster},
		{name: "fastLessAllocs", dedupFunc: DeduplicateLessAllocs},
		{name: "fastLessAllocs2", dedupFunc: DeduplicateLessAllocs2},
		{name: "fastDynamic", dedupFunc: DeduplicateDynamic},
		{name: "fastDynamicLessAllocs", dedupFunc: DeduplicateDynamicLessAllocs},
		{name: "fastDynamicLessAllocs2", dedupFunc: DeduplicateDynamicLessAllocs2},
		{name: "naive", dedupFunc: DeduplicateNaively},
	} {
		t.Run(funcCase.name, func(t *testing.T) {
			for _, tcase := range []struct {
				input    []int
				expected []int
			}{
				{input: nil, expected: nil},
				{input: []int{0}, expected: []int{0}},
				{input: []int{2, 0, 5, 12931293, 10}, expected: []int{0, 2, 5, 10, 12931293}},
				{input: []int{2, 2, 0, 5, 12931293, 5, 10}, expected: []int{0, 2, 5, 10, 12931293}},
				{input: []int{2, 2, 2, 2, 2, 2}, expected: []int{2}},
			} {
				if ok := t.Run("", func(t *testing.T) {
					output := funcCase.dedupFunc(tcase.input)

					// We don't expect any specific order of the output, so sort it on our own.
					sort.Ints(output)
					testutil.Equals(t, tcase.expected, output)
				}); !ok {
					return
				}
			}
		})
	}
}

var output []int

/*
Benchmarks for 1mln elements:

/tmp/GoLand/___BenchmarkDeduplicate_in_github_com_efficientgo_examples_pkg_dups.test -test.v -test.paniconexit0 -test.bench ^\QBenchmarkDeduplicate\E$ -test.run ^$ -test.benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/efficientgo/examples/pkg/dups
cpu: Intel(R) Core(TM) i7-9850H CPU @ 2.60GHz
BenchmarkDeduplicate
BenchmarkDeduplicate/all_duplicates
BenchmarkDeduplicate/all_duplicates/fast
BenchmarkDeduplicate/all_duplicates/fast-12    	    1531	   (3.98s) 3980299 ns/op	       8 B/op	       1 allocs/op
BenchmarkDeduplicate/all_duplicates/fastLessAllocs
BenchmarkDeduplicate/all_duplicates/fastLessAllocs-12         	     585	   9961814 ns/op	22282268 B/op	       2 allocs/op
BenchmarkDeduplicate/all_duplicates/fastLessAllocs2
BenchmarkDeduplicate/all_duplicates/fastLessAllocs2-12        	     460	  12815054 ns/op	22282289 B/op	       3 allocs/op
BenchmarkDeduplicate/all_duplicates/fastDynamic
BenchmarkDeduplicate/all_duplicates/fastDynamic-12            	    2666	   2032222 ns/op	       8 B/op	       1 allocs/op
BenchmarkDeduplicate/all_duplicates/fastDynamicLessAllocs
BenchmarkDeduplicate/all_duplicates/fastDynamicLessAllocs-12  	    2540	   2077542 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeduplicate/all_duplicates/fastDynamicLessAllocs2
BenchmarkDeduplicate/all_duplicates/fastDynamicLessAllocs2-12 	    3046	   1997189 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeduplicate/all_duplicates/naive
BenchmarkDeduplicate/all_duplicates/naive-12                  	    3093	 (2.08ms)  2084348 ns/op	       8 B/op	       1 allocs/op
BenchmarkDeduplicate/no_duplicates
BenchmarkDeduplicate/no_duplicates/fast
BenchmarkDeduplicate/no_duplicates/fast-12                    	      39	 129975298 ns/op	95011827 B/op	   38381 allocs/op
BenchmarkDeduplicate/no_duplicates/fastLessAllocs
BenchmarkDeduplicate/no_duplicates/fastLessAllocs-12          	      57	  92966844 ns/op	22439211 B/op	      20 allocs/op
BenchmarkDeduplicate/no_duplicates/fastLessAllocs2
BenchmarkDeduplicate/no_duplicates/fastLessAllocs2-12         	      66	  82490873 ns/op	22439215 B/op	      20 allocs/op
BenchmarkDeduplicate/no_duplicates/fastDynamic
BenchmarkDeduplicate/no_duplicates/fastDynamic-12             	      44	 132438695 ns/op	95013322 B/op	   38382 allocs/op
BenchmarkDeduplicate/no_duplicates/fastDynamicLessAllocs
BenchmarkDeduplicate/no_duplicates/fastDynamicLessAllocs-12   	      55	  97527610 ns/op	22439256 B/op	      21 allocs/op
BenchmarkDeduplicate/no_duplicates/fastDynamicLessAllocs2
BenchmarkDeduplicate/no_duplicates/fastDynamicLessAllocs2-12  	      62	  85723485 ns/op	22439261 B/op	      21 allocs/op
BenchmarkDeduplicate/no_duplicates/naive
BenchmarkDeduplicate/no_duplicates/naive-12                    	       1	(6.0946657054m) 365679942324 ns/op	(45MB) 45188352 B/op (45MB)	      40 allocs/op (0.365679942 N^2) ns
PASS

Process finished with the exit code 0
*/
func BenchmarkDeduplicate(b *testing.B) {
	for _, tcase := range []struct {
		name  string
		input []int
	}{
		// NOTE: DeduplicateLowAllocs modifies input. However our input slices are
		// prepared for that (one is not modified, due to all dups, second has all zeros so if we
		// returned single element slice with one zero it does not modify array).
		{name: "all duplicates", input: make([]int, 1e6)},
		{name: "no duplicates", input: func() []int {
			input := make([]int, 1e6) // Size: 24B +  1mln * 8B = 8000024B (8MB)
			for i := 0; i < 1e6; i++ {
				input[i] = i
			}
			return input
		}()},
	} {
		b.Run(tcase.name, func(b *testing.B) {
			for _, funcCase := range []struct {
				name      string
				dedupFunc func([]int) []int
			}{
				{name: "fast", dedupFunc: DeduplicateFaster},
				{name: "fastLessAllocs", dedupFunc: DeduplicateLessAllocs},
				{name: "fastLessAllocs2", dedupFunc: DeduplicateLessAllocs2},
				{name: "fastDynamic", dedupFunc: DeduplicateDynamic},
				{name: "fastDynamicLessAllocs", dedupFunc: DeduplicateDynamicLessAllocs},
				{name: "fastDynamicLessAllocs2", dedupFunc: DeduplicateDynamicLessAllocs2},
				//{name: "naive", dedupFunc: DeduplicateNaively}, // Beware, takes ages!
			} {
				b.Run(funcCase.name, func(b *testing.B) {
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						output = funcCase.dedupFunc(tcase.input)
					}
				})
			}
		})
	}
}
