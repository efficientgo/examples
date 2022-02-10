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
	} else {
		end = bytes.LastIndex(b[:end], []byte("\n"))
	}

	// Find last newline before begin and add 1. If not found (-1), it means we
	// are at the start. Otherwise, we start after last newline.
	return bytes.LastIndex(b[:begin], []byte("\n")) + 1, end
}

// ConcurrentSum3 is a basic Sum with added concurrency for introduction
// to go routines. Check ConcurrentSumOpt for the most optimized version.
func ConcurrentSum3(fileName string, workers int) (ret int64, _ error) {
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
			// TODO(bwplotka): Yes, we could optimize bytes.Split a lot, but leaving that for example purposes.
			for _, line := range bytes.Split(b[begin:end], []byte("\n")) {
				// TODO(bwplotka): Yes, we could optimize ParseInt a lot, but leaving that for example purposes.
				num, err := strconv.ParseInt(string(line), 10, 64)
				if err != nil {
					// TODO(bwplotka): Return err using other channel.
					continue
				}
				sum += num
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

// ConcurrentSum2 performs sum concurrently. A lot slower than ConcurrentSum3. An example of pessimisation.
func ConcurrentSum2(fileName string, workers int) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var (
		wg     = sync.WaitGroup{}
		workCh = make(chan []byte, 10)
	)

	wg.Add(workers + 1)
	go func() {
		// TODO(bwplotka): Yes, we could optimize bytes.Split a lot, but leaving that for example purposes.
		for _, line := range bytes.Split(b, []byte("\n")) {
			workCh <- line
		}
		close(workCh)
		wg.Done()
	}()

	for i := 0; i < workers; i++ {
		go func() {
			var sum int64
			for line := range workCh { // Common mistake: for _, line := range <-workCh
				// TODO(bwplotka): Yes, we could optimize ParseInt a lot, but leaving that for example purposes.
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

// ConcurrentSum1 performs sum concurrently. A lot slower than ConcurrentSumOpt. An example of pessimisation.
func ConcurrentSum1(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	for _, line := range bytes.Split(b, []byte("\n")) {
		wg.Add(1)
		// TODO(bwplotka): Yes, we could optimize bytes.Split a lot, but leaving that for example purposes.
		go func(line []byte) {
			defer wg.Done()
			// TODO(bwplotka): Yes, we could optimize ParseInt a lot, but leaving that for example purposes.
			num, err := strconv.ParseInt(string(line), 10, 64)
			if err != nil {
				// TODO(bwplotka): Return err using other channel.
				return
			}
			atomic.AddInt64(&ret, num)
		}(line)
	}

	wg.Wait()
	return ret, nil
}

// Over inline budget, but for readability it's better. Consider splitting functions if needed to get it inlinded.
//./sum_concurrent.go:11:6: cannot inline shardedRange: function too complex: cost 95 exceeds budget 80
func shardedRangeOpt(routineNumber int, bytesPerWorker int, b []byte) (int, int) {
	begin := routineNumber * bytesPerWorker
	end := begin + bytesPerWorker
	if end+bytesPerWorker > len(b) {
		end = len(b)
	}

	// Find last newline before begin and add 1. If not found (-1), it means we
	// are at the start. Otherwise, we start after last newline.
	return bytes.LastIndex(b[:begin], []byte("\n")) + 1, end
}

func ConcurrentSumOpt(fileName string, workers int) (ret int64, _ error) {
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
			begin, end := shardedRangeOpt(i, bytesPerWorker, b)

			var sum int64
			for last := begin; begin < end; begin++ {
				if b[begin] != '\n' {
					continue
				}
				num, err := ParseInt(b[last:begin])
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
