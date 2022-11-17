package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

/*
 command-line-arguments
./interactive2.go:11:30: inlining call to os.Getpid
./interactive2.go:11:30: inlining call to syscall.Getpid
./interactive2.go:11:13: inlining call to fmt.Println
./interactive2.go:16:19: inlining call to os.Open
./interactive2.go:61:13: inlining call to fmt.Println
./interactive2.go:62:12: inlining call to fmt.Scanln
./interactive2.go:64:13: inlining call to fmt.Println
./interactive2.go:66:13: inlining call to fmt.Println
./interactive2.go:67:12: inlining call to fmt.Scanln
./interactive2.go:69:13: inlining call to fmt.Println
./interactive2.go:71:13: inlining call to fmt.Println
./interactive2.go:72:12: inlining call to fmt.Scanln
./interactive2.go:74:13: inlining call to fmt.Println
./interactive2.go:76:13: inlining call to fmt.Println
./interactive2.go:77:12: inlining call to fmt.Scanln
./interactive2.go:79:13: inlining call to fmt.Println
./interactive2.go:80:19: inlining call to os.(*File).Close
./interactive2.go:84:13: inlining call to fmt.Println
./interactive2.go:85:12: inlining call to fmt.Scanln
./interactive2.go:90:13: inlining call to fmt.Println
./interactive2.go:91:12: inlining call to fmt.Scanln
./interactive2.go:97:19: inlining call to os.Open
./interactive2.go:141:13: inlining call to fmt.Println
./interactive2.go:142:13: inlining call to fmt.Println
./interactive2.go:143:13: inlining call to fmt.Println
./interactive2.go:147:13: inlining call to fmt.Println
./interactive2.go:148:19: inlining call to os.(*File).Close
./interactive2.go:11:14: "PID" escapes to heap
./interactive2.go:11:30: int(~R0) escapes to heap
./interactive2.go:11:13: []interface {}{...} does not escape
./interactive2.go:18:12: ... argument does not escape
./interactive2.go:21:11: make([]byte, 600 * 1024 * 1024) escapes to heap
./interactive2.go:24:12: ... argument does not escape
./interactive2.go:27:12: ... argument does not escape
./interactive2.go:27:13: "Read unexpected amount of bytes" escapes to heap
./interactive2.go:27:13: n escapes to heap
./interactive2.go:61:14: "1" escapes to heap
./interactive2.go:61:13: []interface {}{...} does not escape
./interactive2.go:64:14: "Reading 5000 index" escapes to heap
./interactive2.go:64:37: b[5000] escapes to heap
./interactive2.go:64:13: []interface {}{...} does not escape
./interactive2.go:66:14: "2" escapes to heap
./interactive2.go:66:13: []interface {}{...} does not escape
./interactive2.go:69:14: "Reading 100 000 index" escapes to heap
./interactive2.go:69:40: b[100000] escapes to heap
./interactive2.go:69:13: []interface {}{...} does not escape
./interactive2.go:71:14: "3" escapes to heap
./interactive2.go:71:13: []interface {}{...} does not escape
./interactive2.go:74:14: "Reading 104 000 index" escapes to heap
./interactive2.go:74:40: b[104000] escapes to heap
./interactive2.go:74:13: []interface {}{...} does not escape
./interactive2.go:76:14: "4" escapes to heap
./interactive2.go:76:13: []interface {}{...} does not escape
./interactive2.go:79:14: "Close file" escapes to heap
./interactive2.go:79:13: []interface {}{...} does not escape
./interactive2.go:81:12: ... argument does not escape
./interactive2.go:84:14: "Force of memory clear" escapes to heap
./interactive2.go:84:13: []interface {}{...} does not escape
./interactive2.go:90:14: "Finish" escapes to heap
./interactive2.go:90:13: []interface {}{...} does not escape
./interactive2.go:99:12: ... argument does not escape
./interactive2.go:102:11: make([]byte, 600 * 1024 * 1024) escapes to heap
./interactive2.go:108:20: ... argument does not escape
./interactive2.go:108:21: n escapes to heap
./interactive2.go:141:14: "Reading 5000 index" escapes to heap
./interactive2.go:141:37: b[5000] escapes to heap
./interactive2.go:141:13: []interface {}{...} does not escape
./interactive2.go:142:14: "Reading 100 000 index" escapes to heap
./interactive2.go:142:40: b[100000] escapes to heap
./interactive2.go:142:13: []interface {}{...} does not escape
./interactive2.go:143:14: "Reading 104 000 index" escapes to heap
./interactive2.go:143:40: b[104000] escapes to heap
./interactive2.go:143:13: []interface {}{...} does not escape
./interactive2.go:147:14: "Close file" escapes to heap
./interactive2.go:147:13: []interface {}{...} does not escape
./interactive2.go:149:12: ... argument does not escape
<autogenerated>:1: leaking param content: .this
*/

// Buffering 600MB of file in memory.
// Read more in "Efficient Go"; Example 5-2.

func runOpen() {
	fmt.Println("PID", os.Getpid())

	// TODO(bwplotka): Create big file here, so we can play with it - there is no need to upload so big file to GitHub.

	// Open 686MB file and read 600 MB from it.
	f, err := os.Open("test686mbfile.out")
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 600*1024*1024)
	_, err = f.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	// Check out:
	// ps -ax --format=pid,rss,vsz | grep <PID>
	// ls -l /proc/<PID>/map_files
	// cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss

	/*
		c000200000-c025800000 rw-p 00000000 00:00 0 <-- always address of heap
		Size:             612352 kB
		KernelPageSize:        4 kB
		MMUPageSize:           4 kB
		Rss:              612352 kB
		Pss:              612352 kB
		Shared_Clean:          0 kB
		Shared_Dirty:          0 kB
		Private_Clean:         0 kB
		Private_Dirty:    612352 kB
		Referenced:       560216 kB
		Anonymous:        612352 kB
		LazyFree:              0 kB
		AnonHugePages:    352256 kB
		ShmemPmdMapped:        0 kB
		FilePmdMapped:         0 kB
		Shared_Hugetlb:        0 kB
		Private_Hugetlb:       0 kB
		Swap:                  0 kB
		SwapPss:               0 kB
		Locked:                0 kB
		THPeligible:    1
		VmFlags: rd wr mr mw me ac sd hg

	*/

	fmt.Println("1")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 5000 index", b[5000])

	fmt.Println("2")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 100 000 index", b[100000])

	fmt.Println("3")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 104 000 index", b[104000])

	fmt.Println("4")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Close file")
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Force of memory clear")
	fmt.Scanln() // wait for Enter Key

	b = b[0:]
	runtime.GC()

	fmt.Println("Finish")
	fmt.Scanln() // wait for Enter Key

}

func bookExample2() error {
	// Open 686MB file and read 600 MB from it.
	f, err := os.Open("test686mbfile.out")
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 600*1024*1024)
	n, err := f.Read(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("Read unexpected amount of bytes %v", n)
	}

	// Check out:
	// export PID=642103 && ps -ax --format=pid,rss,vsz | grep $PID && cat /proc/$PID/smaps | grep -A22 c000200000-c025800000 | grep Rss
	// cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss
	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss` shows already around 600 MB.
	/*
		c000200000-c025800000 rw-p 00000000 00:00 0 <-- always address of heap
		Size:             612352 kB
		KernelPageSize:        4 kB
		MMUPageSize:           4 kB
		Rss:              612352 kB
		Pss:              612352 kB
		Shared_Clean:          0 kB
		Shared_Dirty:          0 kB
		Private_Clean:         0 kB
		Private_Dirty:    612352 kB
		Referenced:       560216 kB
		Anonymous:        612352 kB
		LazyFree:              0 kB
		AnonHugePages:    352256 kB
		ShmemPmdMapped:        0 kB
		FilePmdMapped:         0 kB
		Shared_Hugetlb:        0 kB
		Private_Hugetlb:       0 kB
		Swap:                  0 kB
		SwapPss:               0 kB
		Locked:                0 kB
		THPeligible:    1
		VmFlags: rd wr mr mw me ac sd hg

	*/

	fmt.Println("Reading 5000 index", b[5000])
	fmt.Println("Reading 100 000 index", b[100000])
	fmt.Println("Reading 104 000 index", b[104000])

	// If we would pause the program in each of those steps `cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss` shows same around 600 MB.

	fmt.Println("Close file")
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss` shows STILL same around 600 MB.

	b = b[0:]
	runtime.GC()

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 c000200000-c025800000 | grep Rss` shows around 140 MB (depends).
	return nil
}
