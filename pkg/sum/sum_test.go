package sum

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/felixge/fgprof"
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
		{f: Sum}, {f: Sum2}, {f: Sum2_scanner}, {f: ConcurrentSum1}, {f: Sum3}, {f: Sum4},
		{f: Sum5}, {f: Sum5_line}, {f: Sum6}, {f: Sum7}, {f: Sum8}, {f: Sum8},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(testFile)
			testutil.Ok(t, err)
			testutil.Equals(t, expectedSum, ret)
		})
	}
}

func TestSumWithWorkers(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	expectedSum, err := createTestInputWithExpectedResult(testFile, 1000)
	testutil.Ok(t, err)

	for _, tcase := range []struct {
		f func(string, int) (int64, error)
	}{
		{f: ConcurrentSum2}, {f: ConcurrentSum3}, {f: ConcurrentSum4},
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

func TestBenchSum(t *testing.T) {
	benchmarkSum(testutil.NewTB(t))
}

func BenchmarkSum_tested(b *testing.B) {
	benchmarkSum(testutil.NewTB(b))
}

func benchmarkSum(tb testutil.TB) {
	fn := lazyCreateTestInput(tb, 2e6)

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

func lazyCreateTestInput(tb testing.TB, numLines int) string {
	tb.Helper()

	fn := fmt.Sprintf("testdata/test.%v.txt", numLines)
	if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
		testutil.Ok(tb, createTestInput(fn, numLines))
	} else {
		testutil.Ok(tb, err)
	}
	return fn
}

// BenchmarkSum benchmarks `Sum` function.
// NOTE(bwplotka): Test it with maximum of 4 CPU cores, given we don't allocate
// more in our production containers.
//
// Recommended run options:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkSum$' -benchtime 5s -count 5 -cpu 1 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkSum(b *testing.B) {
	fn := lazyCreateTestInput(b, 1e8) // 2e6

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sum7(fn)
		testutil.Ok(b, err)
	}
}

// BenchmarkSum_fgprof recommended run options:
// $ export ver=v1fg && go test -run '^$' -bench '^BenchmarkSum_fgprof' -benchtime 60s -count 1 -cpu 1 | tee ${ver}.txt
func BenchmarkSum_fgprof(b *testing.B) {
	f, err := os.Create("fgprof.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFn := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkSum(b)
	testutil.Ok(b, closeFn())
}

func BenchmarkSum_AcrossInputs(b *testing.B) {
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

			fn := lazyCreateTestInput(b, tcase.numLines)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Sum4(fn)
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

var sink string

func BenchmarkStringConv(b *testing.B) {
	b.Run("def", func(b *testing.B) {
		b.ReportAllocs()

		s := make([]byte, 1e6)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sink = string(s)
		}
	})

	b.Run("zero", func(b *testing.B) {
		b.ReportAllocs()

		s := make([]byte, 1e6)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sink = zeroCopyToString(s)
		}
	})

	b.Run("zero?", func(b *testing.B) {
		b.ReportAllocs()

		s := make([]byte, 1e6)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sink = string(zeroCopyToString(s))
		}
	})

	b.Run("copy", func(b *testing.B) {
		b.ReportAllocs()

		s := make([]byte, 1e6)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c := make([]byte, 1e6)
			copy(c, s)
			sink = *((*string)(unsafe.Pointer(&c)))
		}
	})
}
