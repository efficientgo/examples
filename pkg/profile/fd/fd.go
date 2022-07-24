package fd

import (
	"io"
	"os"
	"runtime/pprof"
)

var fdProfile = pprof.NewProfile("github.com/efficient-go/examples/pkg/profile/fd")

// File is a wrapper on os.File that tracks file descriptor lifetime.
type File struct {
	*os.File
}

// Close closes files and updates profile.
func (f *File) Close() error {
	defer fdProfile.Remove(f.File)
	return f.File.Close()
}

// WriteTo writes full profile of the currently open file descriptors in pprof format.
func WriteTo(w io.Writer, debug int) error {
	return fdProfile.WriteTo(w, debug)
}

// Wrap wraps *os.File and to track it in the `fd profile.
// NOTE(bwplotka): We could use finalizers here, but explicit Close is more reliable and accurate.
// Unfortunately it also changes type which might be dropped accidentally.
func Wrap(f *os.File) *File {
	fdProfile.Add(f, 2)
	return &File{File: f}
}

// Open is like os.Open but records `fd` profile.
func Open(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return Wrap(f), nil
}

// Create is like os.Create but records `fd` profile.
func Create(name string) (*File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return Wrap(f), nil
}

// OpenFile is like os.OpenFile but records `fd` profile.
func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return Wrap(f), nil
}
