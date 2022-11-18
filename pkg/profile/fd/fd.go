package fd

import (
	"os"
	"runtime/pprof"
)

// Example custom profile you can write on top `pprof.Profile` helper.
// Read more in "Efficient Go"; Example 9-1.

var fdProfile = pprof.NewProfile("fd.inuse")

// File is a wrapper on os.File that tracks file descriptor lifetime.
type File struct {
	*os.File
}

// Open opens file and tracks it in the `fd` profile`.
// NOTE(bwplotka): We could use finalizers here, but explicit Close is more reliable and accurate.
// Unfortunately it also changes type which might be dropped accidentally.
func Open(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	fdProfile.Add(f, 2)
	return &File{File: f}, nil
}

// Close closes files and updates profile.
func (f *File) Close() error {
	defer fdProfile.Remove(f.File)
	return f.File.Close()
}

// Write saves the profile of the currently open file descriptors in to file in pprof format.
func Write(profileOutPath string) error {
	out, err := os.Create(profileOutPath) // For simplicity, we don't include this file in profile.
	if err != nil {
		return err
	}
	if err := fdProfile.WriteTo(out, 0); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func CreateTemp(dir, pattern string) (*File, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, err
	}
	fdProfile.Add(f, 2)
	return &File{File: f}, nil
}
