package sum

import (
	"errors"
	"fmt"
	"os"
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

func createTestInput(in string) error {
	fmt.Println("create")
	// YOLO
	return nil
}

// BenchmarkSum benchmarks `Sum` function.
// NOTE(bwplotka): Test it with maximum of 4 CPU cores, given we don't allocate
// more in our production containers.
//
// Recommended run options:
// $ export var=v1 && go test -count 5 -benchmem -benchtime 5s -run '^$'  -bench '^BenchmarkSum$' \
// -cpu=4 -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt
func BenchmarkSum(b *testing.B) {
	// Ensure test input was generated.
	// NOTE(bwplotka): Change name of test input if you choose to change the test data.
	if _, err := os.Stat("testdata/test.2M.txt"); errors.Is(err, os.ErrNotExist) {
		testutil.Ok(b, createTestInput("testdata/test.2M.txt"))
	} else {
		testutil.Ok(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sum("testdata/test.1M.txt")
		testutil.Ok(b, err)
	}
}

// BenchmarkSum
//BenchmarkSum-12    	      15	  74969706 ns/op	60917664 B/op	 1636366 allocs/op
// BenchmarkSum-12    	      14	  75451998 ns/op	60917678 B/op	 1636366 allocs/op
