package sum

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/pkg/errors"
)

const testSumOfTen = int64(31108)

func TestSum(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	testutil.Ok(t, createTestInput(testFile, 1000))

	for _, tcase := range []struct {
		f func(string) (int64, error)
	}{
		{f: Sum}, {f: Sum2}, {f: ConcurrentSum1}, {f: Sum3}, {f: Sum4},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(testFile)
			testutil.Ok(t, err)
			testutil.Equals(t, 100*testSumOfTen, ret)
		})
	}
}

func TestSumWithWorkers(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "input.txt")
	testutil.Ok(t, createTestInput(testFile, 1000))

	for _, tcase := range []struct {
		f func(string, int) (int64, error)
	}{
		{f: ConcurrentSum2}, {f: ConcurrentSum3}, {f: ConcurrentSumOpt},
	} {
		t.Run("", func(t *testing.T) {
			ret, err := tcase.f(testFile, 4)
			testutil.Ok(t, err)
			testutil.Equals(t, 100*testSumOfTen, ret)

			ret, err = tcase.f(testFile, 11)
			testutil.Ok(t, err)
			testutil.Equals(t, 100*testSumOfTen, ret)
		})
	}
}

func createTestInput(in string, numLen int) (err error) {
	if numLen%10 != 0 {
		return errors.Errorf("number of input should be division by 10, got %v", numLen)
	}

	if err := os.MkdirAll(filepath.Dir(in), os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(in, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "open")
	}

	for i := 0; i < numLen/10; i++ {
		if _, err := f.WriteString(`123
43
632
22
2
122
26660
91
2
3411
`); err != nil {
			return err
		}
	}

	return f.Close()
}

// BenchmarkSum benchmarks `Sum` function.
// NOTE(bwplotka): Test it with maximum of 4 CPU cores, given we don't allocate
// more in our production containers.
//
// Recommended run options:
// $ export ver=v1-2 && go test -count 5 -benchmem -benchtime 10s -run '^$' -bench '^BenchmarkSum$' -cpu 4 -memprofile ${ver}.mem.pprof -cpuprofile ${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkSum(b *testing.B) {
	// Ensure test input was generated.
	// NOTE(bwplotka): Change name of test input if you choose to change the test data.
	if _, err := os.Stat("testdata/test.2M.txt"); errors.Is(err, os.ErrNotExist) {
		testutil.Ok(b, createTestInput("testdata/test.2M.txt", 2e6)) // ~7.55 MB, 2 million lines.
	} else {
		testutil.Ok(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConcurrentSum3("testdata/test.2M.txt", 4)
		testutil.Ok(b, err)
	}
}
