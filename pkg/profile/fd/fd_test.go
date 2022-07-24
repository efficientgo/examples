package fd_test

import (
	"io"
	"os"
	"sync"
	"testing"

	"github.com/efficientgo/examples/pkg/profile/fd"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

type testApp struct {
	files []io.ReadCloser
}

func (a *testApp) Close() {
	for _, cl := range a.files {
		_ = cl.Close()
	}
}

func (a *testApp) Open(fName string) error {
	f, err := os.Open(fName)
	if err != nil {
		return err
	}

	a.files = append(a.files, fd.Wrap(f))
	return nil
}

//go:noinline
func (a *testApp) funcD(file string) {
	for i := 0; i < 10; i++ {
		_ = a.Open(file) // TODO: Report error...
	}
}

//go:noinline
func (a *testApp) funcC(file string) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			a.funcD(file)
		}
		wg.Done()
	}()
	wg.Wait()
}

//go:noinline
func (a *testApp) funcB(file string) {
	a.funcD(file)
}

//go:noinline
func (a *testApp) funcA(file string) {
	_ = a.Open(file) // TODO: Report error...
}

func TestFD(t *testing.T) {
	a := &testApp{}
	t.Cleanup(a.Close)

	a.funcA("/dev/null")
	a.funcB("/dev/null")
	a.funcC("/dev/null")

	out, err := os.Create("fd.pprof")
	testutil.Ok(t, err)
	testutil.Ok(t, fd.WriteTo(out, 0))
	testutil.Ok(t, out.Close())
}
