package sum

import (
	"runtime"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

const (
	inputFileName = "input.txt"
	inputSum      = int64(242028430)
)

func TestSum(t *testing.T) {
	for _, tcase := range []struct {
		f func(string) (int64, error)
	}{
		{f: Sum}, {f: Sum2}, {f: ConcurrentSum1}, {f: Sum3}, {f: Sum4},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(inputFileName) // 3.55 MB 1mln lines
			testutil.Ok(t, err)
			testutil.Equals(t, inputSum, ret)
		})
	}
}

func TestSumWithWorkers(t *testing.T) {
	for _, tcase := range []struct {
		f func(string, int) (int64, error)
	}{
		{f: ConcurrentSum2}, {f: ConcurrentSum3}, {f: ConcurrentSumOpt},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(inputFileName, 4) // 3.55 MB 1mln lines
			testutil.Ok(t, err)
			testutil.Equals(t, inputSum, ret)

			ret, err = tcase.f(inputFileName, 11) // 3.55 MB 1mln lines
			testutil.Ok(t, err)
			testutil.Equals(t, inputSum, ret)
		})
	}
}

var Answer int64

// export var=v1 && go test -count 5 -benchtime 5s -run '^$' -bench . -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt
func BenchmarkSum(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer, _ = Sum3("input2.txt")
	}
}
