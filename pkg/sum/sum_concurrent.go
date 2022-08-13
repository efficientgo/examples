package sum

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"

	"github.com/efficientgo/tools/core/pkg/errcapture"
)

// ConcurrentSum1 performs sum concurrently. A lot slower than ConcurrentSum3. An example of pessimisation.
func ConcurrentSum1(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	var last int
	for curr := 0; curr < len(b); curr++ {
		if b[curr] != '\n' {
			continue
		}

		wg.Add(1)
		go func(line []byte) { // Creation of goroutine turns to be mem intensive on scale! (on top of time)
			defer wg.Done()
			num, err := ParseInt(line)
			if err != nil {
				// TODO(bwplotka): Return err using other channel.
				return
			}
			atomic.AddInt64(&ret, num)
		}(b[last:curr])
		last = curr + 1
	}
	wg.Wait()
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
		var last int
		for curr := 0; curr < len(b); curr++ {
			if b[curr] != '\n' {
				continue
			}
			workCh <- b[last:curr]
			last = curr + 1
		}
		close(workCh)
		wg.Done()
	}()

	for i := 0; i < workers; i++ {
		go func() {
			var sum int64
			for line := range workCh { // Common mistake: for _, line := range <-workCh
				num, err := ParseInt(line)
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

func ConcurrentSum4(fileName string, workers int) (ret int64, _ error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	s, err := f.Stat()
	if err != nil {
		return 0, err
	}

	var (
		bytesPerWorker = int(s.Size()) / workers
		resultCh       = make(chan int64)
	)

	for i := 0; i < workers; i++ {
		go func(i int) {
			buf := bytes.Buffer{}
			buf.Grow(bytesPerWorker + 10) // Assuming 10 max digit

			// Coordination-free algorithm, which shards buffered file deterministically.
			begin := i * bytesPerWorker
			end := begin + bytesPerWorker
			if end+bytesPerWorker > int(s.Size()) {
				end = int(s.Size())
			}

			_, err := io.Copy(&buf, io.NewSectionReader(f, int64(begin), int64(bytesPerWorker+10)))
			if err != nil {
				// TODO(bwplotka): Return err using other channel.
				fmt.Println(err)
				resultCh <- 0
				return
			}

			var (
				last int
				sum  int64
				b    = buf.Bytes()
			)
			for i := range b {
				if b[i] != '\n' {
					continue
				}
				if last > 0 || begin == 0 {
					num, err := ParseInt(b[last:i])
					if err != nil {
						// TODO(bwplotka): Return err using other channel.
						fmt.Println(err)
						continue
					}
					sum += num
				}
				last = i + 1

				if begin+last > end {
					break
				}
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

func shardedEstimatedRange(routineNumber int, bytesPerWorker int, size int) (int, int) {
	begin := routineNumber * bytesPerWorker
	end := begin + bytesPerWorker
	if end+bytesPerWorker > size {
		end = size
	}
	return begin, end
}

func ConcurrentSum4_buf(fileName string, workers int) (ret int64, _ error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}

	s, err := f.Stat()
	if err != nil {
		return 0, err
	}

	var (
		bytesPerWorker = int(s.Size()) / workers
		resultCh       = make(chan int64)
	)

	for i := 0; i < workers; i++ {
		go func(i int) {
			buf := make([]byte, 512*1024)

			// Coordination-free algorithm, which shards buffered file deterministically.
			begin, end := shardedEstimatedRange(i, bytesPerWorker, int(s.Size()))
			r := io.NewSectionReader(f, int64(begin), s.Size())

			var (
				readOff int
				oneMore bool
				sum     int64
				err     error
				n       int
			)

		bigLoop:
			for err != io.EOF {
				n, err = r.ReadAt(buf, int64(readOff))
				if err != nil && err != io.EOF {
					// TODO(bwplotka): Return err using other channel.
					fmt.Println(err)
					break
				}

				var last int
				for i := range buf[:n] {
					if buf[i] != '\n' {
						continue
					}
					if last > 0 || begin == 0 {
						num, err := ParseInt(buf[last:i])
						if err != nil {
							// TODO(bwplotka): Return err using other channel.
							fmt.Println(err)
							continue
						}
						sum += num
						last = i + 1
					}
					readOff += last

					if begin+readOff > end {
						if oneMore {
							break bigLoop
						}
						oneMore = true
					}
				}
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
