package leak

import (
	"io"
	"os"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/logerrcapture"
	"github.com/efficientgo/core/merrors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Example on common leaks using `os.File`.
// Read more in "Efficient Go"; Example 11-8.

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

// Example on common leaks using `os.File` when multiple files are used.
// Read more in "Efficient Go"; Example 11-9.

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

func closeAll(closers []io.ReadCloser) error {
	errs := merrors.New()
	for _, c := range closers {
		errs.Add(c.Close())
	}
	return errs.Err()
}
