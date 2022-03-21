package main

import (
	"fmt"
	"log"
	"os"

	"github.com/efficientgo/examples/pkg/memory/mmap"
)

func main() {
	fmt.Println("PID", os.Getpid())

	// TODO(bwplotka): Create big file here, so we can play with it - there is no need to upload so big file to GH.

	// Mmap 600 MB of 686MB file.
	f, err := mmap.OpenFile("test686mbfile.out", 600*1024*1024)
	if err != nil {
		log.Fatal(err)
	}

	// Check out:
	// ls -l /proc/<PID>/map_files
	// cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss
	b := f.Bytes()

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

	fmt.Println("Unmapping")
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finish")
	fmt.Scanln() // wait for Enter Key

}

func bookExample1() error {
	// Mmap 600MB of 686MB file.
	f, err := mmap.OpenFile("test686mbfile.out", 600*1024*1024)
	if err != nil {
		return err
	}
	b := f.Bytes()

	// At this point we can see symlink to test686mbfile.out file in /proc/<PID>/map_files.
	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 0KB.
	fmt.Println("Reading 5000 index", b[5000])

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 48-70KB.

	fmt.Println("Reading 100 000 index", b[100000])

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 100-126KB.

	fmt.Println("Reading 104 000 index", b[104000])

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 100-126KB (same).

	if err := f.Close(); err != nil {
		return err
	}

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows nothing, RSS freed.
	return nil
}
