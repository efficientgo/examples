package leak

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/efficientgo/tools/core/pkg/logerrcapture"
	"github.com/efficientgo/tools/core/pkg/merrors"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func doWithFile_Wrong(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close() // Wrong!

	// Use file...

	return nil
}

func doWithFile_LogCloseErr(logger log.Logger, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		level.Error(logger).Log("err", err)
		return
	}
	defer logerrcapture.Do(logger, f.Close, "close file")

	// Use file...
}

func doWithFile_CaptureCloseErr(fileName string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	// Use file...

	return nil
}

func TestDoWithFile(t *testing.T) {
	testutil.Ok(t, doWithFile_Wrong("/dev/null"))
	testutil.Ok(t, doWithFile_CaptureCloseErr("/dev/null"))
	doWithFile_LogCloseErr(log.NewNopLogger(), "/dev/null")
}

func openMultiple_Wrong(fileNames ...string) ([]io.ReadCloser, error) {
	files := make([]io.ReadCloser, 0, len(fileNames))
	for _, fn := range fileNames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, err // Leaked files!
		}
		files = append(files, f)
	}
	return files, nil
}

func closeAll(closers []io.ReadCloser) error {
	errs := merrors.New()
	for _, c := range closers {
		errs.Add(c.Close())
	}
	return errs.Err()
}

func openMultiple_Correct(fileNames ...string) ([]io.ReadCloser, error) {
	files := make([]io.ReadCloser, 0, len(fileNames))
	for _, fn := range fileNames {
		f, err := os.Open(fn)
		if err != nil {
			return nil, merrors.New(err, closeAll(files)).Err()
		}
		files = append(files, f)
	}
	return files, nil
}

func TestOpenMultiple(t *testing.T) {
	files, err := openMultiple_Wrong("/dev/null", "/dev/null", "/dev/null")
	testutil.Ok(t, err)
	testutil.Ok(t, closeAll(files))

	files, err = openMultiple_Correct("/dev/null", "/dev/null", "/dev/null")
	testutil.Ok(t, err)
	testutil.Ok(t, closeAll(files))
}

func AmazingConcurrentCode() {
	time.Sleep(199 * time.Millisecond)
}

func BenchmarkAmazingConcurrentCode_Wrong(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() { AmazingConcurrentCode() }()
		go func() { AmazingConcurrentCode() }()
	}
}

func BenchmarkAmazingConcurrentCode_Better(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			AmazingConcurrentCode()
		}()
		go func() {
			defer wg.Done()
			AmazingConcurrentCode()
		}()

		wg.Wait()
	}
}
