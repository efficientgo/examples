package sum

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"sync/atomic"
)

// Over inline budget, but for readability it's better. Consider splitting functions if needed to get it inlinded.
//./sum_concurrent.go:11:6: cannot inline shardedRange: function too complex: cost 95 exceeds budget 80
func shardedRange(routineNumber int, bytesPerWorker int, b []byte) (int, int) {
	begin := routineNumber * bytesPerWorker
	end := begin + bytesPerWorker
	if end+bytesPerWorker > len(b) {
		end = len(b)
	}

	// Find last newline before begin and add 1. If not found (-1), it means we
	// are at the start. Otherwise, we start after last newline.
	return bytes.LastIndex(b[:begin], []byte("\n")) + 1, end
}

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
		go func(i int) {
			// Coordination-free algorithm, which shards buffered file deterministically.
			begin, end := shardedRange(i, bytesPerWorker, b)

			var sum int64
			for last := begin; begin < end; begin++ {
				if b[begin] != '\n' {
					continue
				}
				num, err := strconv.ParseInt(string(b[last:begin]), 10, 64)
				if err != nil {
					// TODO(bwplotka): Return err using other channel.
					continue
				}
				sum += num
				last = begin + 1
			}
			resultCh <- sum
		}(i)
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
