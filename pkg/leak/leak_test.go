package leak

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/logerrcapture"
	"github.com/efficientgo/core/merrors"
	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.uber.org/goleak"
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

func ComplexComputation() int {
	time.Sleep(1 * time.Second) // Computation.
	time.Sleep(1 * time.Second) // Cleanup.
	return 4
}

func Handle_VeryWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

func Handle_Wrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

func ComplexComputationWithCtx(ctx context.Context) (ret int) {
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second): // Computation.
		ret = 4
	}

	time.Sleep(1 * time.Second) // Cleanup.
	return ret
}

func Handle_AlsoWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context())
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

func Handle_Better(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context())
	}()

	resp := <-respCh
	if r.Context().Err() != nil {
		return
	}

	_, _ = w.Write([]byte(strconv.Itoa(resp)))
}

func TestHandle(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_VeryWrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_Wrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_AlsoWrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_Better(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

}

func TestHandleCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "https://efficientgo.com", nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		Handle_VeryWrong(w, r.WithContext(ctx))
		wg.Done()
	}()
	// Immediately cancel.
	cancel()

	time.Sleep(3 * time.Second)
	wg.Wait()
}

func BenchmarkComplexComputation_Better(b *testing.B) {
	defer goleak.VerifyNone(
		b,
		goleak.IgnoreTopFunction("testing.(*B).run1"),
		goleak.IgnoreTopFunction("testing.(*B).doBench"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			ComplexComputation()
		}()
		go func() {
			defer wg.Done()
			ComplexComputation()
		}()

		wg.Wait()
	}
}
