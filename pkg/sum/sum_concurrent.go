package sum

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"sync/atomic"
)

func ConcurrentSum(fileName string, workers int) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var (
		bytesPerWorker = len(b) / workers
		resultCh       = make(chan int64)
	)

	for i := 0; i < workers; i++ {
		// Coordination-free algorithm, which shards and reuses buffered file deterministically.
		go func(begin int) {
			end := begin + bytesPerWorker
			if begin+2*bytesPerWorker > len(b) {
				end = len(b)
			} else {
				end = bytes.LastIndex(b[begin:end], []byte("\n"))
			}

			// Find last newline before begin and add 1. If not found (-1), it means we
			// are at the start. Otherwise we start after last newline.
			begin = bytes.LastIndex(b[:begin], []byte("\n")) + 1

			var sum int64
			for _, line := range bytes.Split(b[begin:end], []byte("\n")) {
				num, err := strconv.ParseInt(string(line), 10, 64)
				if err != nil {
					// TODO(bwplotka): Return err using other channel.
					continue
				}
				sum += num
			}
			resultCh <- sum
		}(i * bytesPerWorker)
	}

	for i := 0; i < workers; i++ {
		ret += <-resultCh
	}
	close(resultCh)
	return ret, nil
}

// ConcurrentSum2 performs sum concurrently. A lot slower than ConcurrentSum. An example of pessimisation.
func ConcurrentSum2(fileName string, workers int) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var (
		wg     = sync.WaitGroup{}
		workCh = make(chan []byte, 10)
		last   int
	)

	wg.Add(workers + 1)
	go func() {
		// TODO(bwplotka): Stream it from file.
		for i := 0; i < len(b); i++ {
			if b[i] == '\n' {
				workCh <- b[last:i]
				last = i + 1
			}
		}
		close(workCh)
		wg.Done()
	}()
	for i := 0; i < workers; i++ {
		go func() {
			var sum int64

			for line := range workCh { // Common mistake: for _, line := range <-workCh
				num, err := strconv.ParseInt(string(line), 10, 64)
				if err != nil {
					// TODO(bwplotka): Return err using other channel.
					continue
				}
				sum += num
			}
			atomic.AddInt64(&ret, sum)
			wg.Done()
		}()
	}
	wg.Wait()
	return ret, nil
}
