package fd

import (
	"os"
	"runtime/pprof"
)

var fdProfile = pprof.NewProfile("fd.inuse")

// Wrap wraps *os.File to track it in the `fd` profile`.
// NOTE(bwplotka): We could use finalizers here, but explicit Close is more reliable and accurate.
// Unfortunately it also changes type which might be dropped accidentally.
func Wrap(f *os.File) *File {
	fdProfile.Add(f, 2)
	return &File{File: f}
}

// File is a wrapper on os.File that tracks file descriptor lifetime.
type File struct {
	*os.File
}

// Close closes files and updates profile.
func (f *File) Close() error {
	defer fdProfile.Remove(f.File)
	return f.File.Close()
}

// Write saves the profile of the currently open file descriptors in to file in pprof format.
func Write(profileOutPath string) error {
	out, err := os.Create(profileOutPath)
	if err != nil {
		return err
	}
	if err := fdProfile.WriteTo(out, 0); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
