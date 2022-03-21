package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

var m = &runtime.MemStats{}

func createSlice(size int) {
	b := make([]byte, size)
	{
		runtime.GC()
		runtime.ReadMemStats(m)
		fmt.Println(m.HeapAlloc)
		fmt.Scanln() // wait for Enter Key
	}

	for i := range b {
		b[i] = 'a'
	}

	{
		runtime.GC()
		runtime.ReadMemStats(m)
		fmt.Println(m.HeapAlloc)
		fmt.Scanln() // wait for Enter Key
	}
}

func main() {
	{
		runtime.GC()
		runtime.ReadMemStats(m)
		fmt.Println(m.HeapAlloc)
	}

	createSlice(4097)
	fmt.Println(unsafe.Sizeof(yolo))
}
