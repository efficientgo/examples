package sum

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/pkg/errors"
)

func createTestInput(fn string, numLen int) error {
	_, err := createTestInputWithExpectedResult(fn, numLen)
	return err
}

func createTestInputWithExpectedResult(fn string, numLen int) (sum int64, err error) {
	if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
		return 0, err
	}

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, errors.Wrap(err, "open")
	}

	defer errcapture.Do(&err, f.Close, "close file")

	return sumtestutil.CreateTestInputWithExpectedResult(f, numLen)
}

func TestSum(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	expectedSum, err := createTestInputWithExpectedResult(testFile, 1000)
	testutil.Ok(t, err)

	for _, tcase := range []struct {
		f func(string) (int64, error)
	}{
		{f: Sum}, {f: Sum2}, {f: ConcurrentSum1}, {f: Sum3}, {f: Sum4},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(testFile)
			testutil.Ok(t, err)
			testutil.Equals(t, expectedSum, ret)
		})
	}
}

func TestSumSimple(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	expectedSum, err := createTestInputWithExpectedResult(testFile, 1000)
	testutil.Ok(t, err)

	ret, err := Sum(testFile)
	testutil.Ok(t, err)
	testutil.Equals(t, expectedSum, ret)
}

func TestSumWithWorkers(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	expectedSum, err := createTestInputWithExpectedResult(testFile, 1000)
	testutil.Ok(t, err)

	for _, tcase := range []struct {
		f func(string, int) (int64, error)
	}{
		{f: ConcurrentSum2}, {f: ConcurrentSum3}, {f: ConcurrentSumOpt},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(testFile, 4)
			testutil.Ok(t, err)
			testutil.Equals(t, expectedSum, ret)

			ret, err = tcase.f(testFile, 11)
			testutil.Ok(t, err)
			testutil.Equals(t, expectedSum, ret)
		})
	}
}

// BenchmarkSum benchmarks `Sum` function.
// NOTE(bwplotka): Test it with maximum of 4 CPU cores, given we don't allocate
// more in our production containers.
//
// Recommended run options:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkSum1$' -benchtime 5s -count 5 -cpu 4 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkSum1(b *testing.B) {
	// Create ~7.55 MB file with 2 million lines.
	fn := filepath.Join(b.TempDir(), "/test.2M.txt")
	testutil.Ok(b, createTestInput(fn, 2e6))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sum(fn)
		testutil.Ok(b, err)
	}
}

func BenchmarkSum(b *testing.B) {
	for _, tcase := range []struct {
		numLines int
	}{
		{numLines: 0},
		{numLines: 1e2},
		{numLines: 1e4},
		{numLines: 1e6},
		{numLines: 2e6},
	} {
		b.Run(fmt.Sprintf("lines-%d", tcase.numLines), func(b *testing.B) {
			b.ReportAllocs()

			fn := fmt.Sprintf("testdata/test.%v.txt", tcase.numLines)
			testutil.Ok(b, createTestInput(fn, tcase.numLines))

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Sum(fn)
				testutil.Ok(b, err)
			}
		})
	}
}

func BenchmarkAssert(b *testing.B) {
	b.ReportAllocs()

	var err error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testutil.Ok(b, err)
	}
}

func TestBenchSum2(t *testing.T) {
	benchmarkSum(testutil.NewTB(t))
}

func BenchmarkSum2(b *testing.B) {
	benchmarkSum(testutil.NewTB(b))
}

func benchmarkSum(tb testutil.TB) {
	const numLines int = 2e6 // ~7.55 MB, 2 million lines.

	fn := fmt.Sprintf("testdata/test.%v.txt", numLines)
	if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
		testutil.Ok(tb, createTestInput(fn, numLines))
	} else {
		testutil.Ok(tb, err)
	}

	tb.ResetTimer()
	for i := 0; i < tb.N(); i++ {
		ret, err := Sum(fn)
		testutil.Ok(tb, err)
		if !tb.IsBenchmark() {
			// More expensive result checks can be here.
			testutil.Equals(tb, int64(6221600000), ret)
		}
	}
}

func BenchmarkSum0(b *testing.B) {
	runtime.MemProfileRate = 1 // Not advised on production!

	for i := 0; i < b.N; i++ {
		_, err := Sum("./testdata/test.2000000.txt")
		testutil.Ok(b, err)
	}
}
