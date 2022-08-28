package file

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
)

func createTestInput(fn string, bytes int) (err error) {
	if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "open")
	}

	defer func() {
		if err != nil {
			errcapture.Do(&err, func() error { return os.Remove(fn) }, "remove failed file")
		}
	}()

	b := make([]byte, 10*1024*1024)
	for i := 0; i < bytes/len(b); i++ {

		if _, err := rand.Read(b); err != nil {
			return errors.Wrap(err, "read urandom")
		}
		if _, err := f.Write(b); err != nil {
			return err
		}
	}

	return f.Close()
}

const partSize int64 = 1 * 1024 * 1024

func readThreeSections(t *testing.T, f *os.File) ([]byte, []byte, []byte) {
	s1 := io.NewSectionReader(f, 20, partSize)
	s2 := io.NewSectionReader(f, 2*partSize, partSize)
	s3 := io.NewSectionReader(f, 10*partSize, partSize)

	s1read := make([]byte, partSize)
	_, err := s1.Read(s1read)
	testutil.Ok(t, err)

	s3read := make([]byte, partSize)
	_, err = s3.Read(s3read)
	testutil.Ok(t, err)

	s2read := make([]byte, partSize)
	_, err = s2.Read(s2read)
	testutil.Ok(t, err)

	return s1read, s2read, s3read
}

func TestFile_SectionReader(t *testing.T) {
	fn := filepath.Join(t.TempDir(), "test.file")
	testutil.Ok(t, createTestInput(fn, 100*1024*1024))

	f, err := os.Open(fn)
	testutil.Ok(t, err)

	s1read, s2read, s3read := readThreeSections(t, f)

	testutil.Assert(t, bytes.Compare(s1read, s2read) != 0)
	testutil.Assert(t, bytes.Compare(s2read, s3read) != 0)
	testutil.Assert(t, bytes.Compare(s1read, s3read) != 0)

	defer func() { testutil.Ok(t, f.Close()) }()

	// We can reuse same file descriptor for new section reads - they are not touching "current offset".

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			s1, s2, s3 := readThreeSections(t, f)
			testutil.Equals(t, s1read, s1)
			testutil.Equals(t, s2read, s2)
			testutil.Equals(t, s3read, s3)
		}()
	}
	wg.Wait()
}
