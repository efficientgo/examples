package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/efficientgo/examples/pkg/profile/fd"
)

type TestApp struct {
	files []io.ReadCloser
}

func (a *TestApp) Close() {
	for _, cl := range a.files {
		_ = cl.Close() // TODO: Check error.
	}
	a.files = a.files[:0]
}

func (a *TestApp) open(fName string) {
	f, _ := os.Open(fName) // TODO: Check error.
	a.files = append(a.files, fd.Wrap(f))
}

func (a *TestApp) OpenSingleFile(file string) {
	a.open(file)
}

func (a *TestApp) OpenTenFiles(file string) {
	for i := 0; i < 10; i++ {
		a.open(file)
	}
}

func (a *TestApp) Open100FilesConcurrently(file string) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			a.OpenTenFiles(file)
			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {
	a := &TestApp{}
	defer a.Close()

	// No matter how many files we opened in the past...
	for i := 0; i < 10; i++ {
		a.OpenTenFiles("/dev/null")
		a.Close()
	}

	// ...after last close only currently used will be in profile.
	a.OpenSingleFile("/dev/null")
	a.OpenTenFiles("/dev/null")
	a.Open100FilesConcurrently("/dev/null")

	if err := fd.Write("fd.pprof"); err != nil {
		log.Fatal(err)
	}
}
