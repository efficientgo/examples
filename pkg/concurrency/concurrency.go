package concurrency

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

func anotherFunction(arg1 string) { /*...*/ }

func function() {
	// Scope of the current goroutine.
	// ...

	go func() {
		// This scope will run concurrently any moment now.
		// ...
	}()

	// anotherFunction will run concurrently any moment now.
	go anotherFunction("argument1")

	// After our function ends, two goroutines we started can still run.
	return
}

var randInt64 = func() int64 {
	return rand.Int63()
}

func sharingWithSeparateSpace() (sum int64) {
	var wg sync.WaitGroup
	results := [3]int64{}

	concurrentFn := func(i int) {
		// ...
		results[i] = randInt64()
		wg.Done()
	}
	wg.Add(3)
	go concurrentFn(0)
	go concurrentFn(1)
	go concurrentFn(2)

	wg.Wait()
	for _, res := range results {
		sum += res
	}
	return sum
}

func sharingWithAtomic() (sum int64) {
	var wg sync.WaitGroup

	concurrentFn := func() {
		// ...
		atomic.AddInt64(&sum, randInt64())
		wg.Done()
	}
	wg.Add(3)
	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	wg.Wait()
	return sum
}

func sharingWithMutex() (sum int64) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	concurrentFn := func() {
		// ...
		mu.Lock()
		sum += randInt64()
		mu.Unlock()
		wg.Done()
	}
	wg.Add(3)
	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	wg.Wait()
	return sum
}

func sharingWithChannel() (sum int64) {
	result := make(chan int64)

	concurrentFn := func() {
		// ...
		result <- randInt64()
	}
	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	for i := 0; i < 3; i++ {
		sum += <-result
	}
	close(result)
	return sum
}

func sharingWithChannel2() (sum int64) {
	result := make(chan int64)

	concurrentFn := func() {
		// ...
		result <- randInt64()
		result <- randInt64()
		result <- randInt64()
		close(result)
	}
	go concurrentFn()

	for res := range result {
		sum += res
	}
	return sum
}
