package sum

import (
	"bytes"
	"fmt"
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
		lines          = bytes.Split(b, []byte("\n"))
		linesPerWorker = len(lines) / workers
		resultCh       = make(chan int64)
	)

	for i := 0; i < workers; i++ {
		end := (i + 1) * linesPerWorker
		if i == workers-1 {
			end = len(lines) // Last one takes all.
		}
		go func(begin, end int) {
			var sum int64

			for _, line := range lines[begin:end] {
				num, err := strconv.ParseInt(string(line), 10, 64)
				if err != nil {
					// TODO(bwplotka): Return err using other channel.
					continue
				}
				sum += num
			}
			resultCh <- sum
		}(i*linesPerWorker, end)
	}

	for i := 0; i < workers; i++ {
		ret += <-resultCh
	}
	close(resultCh)
	return ret, nil
}

func ConcurrentSum2(fileName string, workers int) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var (
		wg     = sync.WaitGroup{}
		workCh = make(chan []byte)
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
		fmt.Println("done")
		close(workCh)
		wg.Done()
	}()
	for i := 0; i < workers; i++ {
		go func() {
			var sum int64

			for _, line := range <-workCh {
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
