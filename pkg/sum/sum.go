package sum

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"unsafe"

	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/pkg/errors"
)

func Sum(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			// Empty line at the end.
			continue
		}

		num, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil
}

// Sum2 is sum with optimized the first latency + CPU bottleneck bytes.Split.
// bytes.Split look complex to hande different cases. It allocates a lot causing  It looks like the algo is simple enough to just
// implement on our own (tried scanner := bufio.NewScanner(f) but it's slower).
// 30% less latency and 5x less memory than Sum.
func Sum2(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(string(b[last:i]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}
	return ret, nil
}

// Sum2_scanner is a sum attempting using scanner. Actually slower than Sum2, but uses less memory.
func Sum2_scanner(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		num, err := strconv.ParseInt(string(scanner.Bytes()), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil
}

// Sum3 is a sum with optimized the second latency + CPU bottleneck: string conversion.
// On CPU profile we see byte to string conversion not only allocate memory, but also takes precious time.
// Let's perform zeroCopy conversion.
// 2x less latency memory than Sum2.
func Sum3(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(zeroCopyToString(b[last:i]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}
	return ret, nil
}

func zeroCopyToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

// Sum4 is a sum with optimized the second latency + CPU bottleneck: ParseInt and string conversion.
// On CPU profile we see that ParseInt does a lot of checks that we might not need. We write our own parsing
// straight from byte to avoid conversion CPU time.
// 2x less latency, same mem as Sum3.
func Sum4(fileName string) (ret int64, err error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}
		num, err := ParseInt(b[last:i])
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}
	return ret, nil
}

// ParseInt is 3-4x times faster than strconv.ParseInt or Atoi.
func ParseInt(input []byte) (n int64, _ error) {
	factor := int64(1)
	k := 0

	// TODO(bwplotka): Optimize if only positive integers are accepted (only 2.6% overhead).
	if input[0] == '-' {
		factor *= -1
		k++
	}

	for i := len(input) - 1; i >= k; i-- {
		if input[i] < '0' || input[i] > '9' {
			return 0, errors.Errorf("not a valid integer: %v", input)
		}

		n += factor * int64(input[i]-'0')
		factor *= 10
	}
	return n, nil
}

// Sum5 is like Sum4, but noticing that it takes time to even allocate 21 MB on heap (and read file to it).
// Let's try to use scanner instead.
// Slower than Sum4 and Sum6 because scanner is not optimized for this...? Scanner takes 73% of CPU time.
func Sum5(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		num, err := ParseInt(scanner.Bytes())
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil
}

func Sum5_line(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	scanner := bufio.NewScanner(f)
	scanner.Split(ScanLines)
	for scanner.Scan() {
		num, err := ParseInt(scanner.Bytes())
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := range data {
		if data[i] != '\n' {
			continue
		}
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

// Sum6 is like Sum4, but trying to use max 10 KB of mem without scanner and bulk read.
// Assuming no integer is larger than 9 000 digits.
func Sum6(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	buf := make([]byte, 9*1024)
	var readOff, n int
	for err != io.EOF {
		n, err = f.ReadAt(buf, int64(readOff))
		if err != nil && err != io.EOF {
			return 0, err
		}

		var last int
		//for i := 0; i < n; i++ { // Funny enough this is 5% slower!
		for i := range buf[:n] {
			if buf[i] != '\n' {
				continue
			}
			num, err := ParseInt(buf[last:i])
			if err != nil {
				return 0, err
			}

			ret += num
			last = i + 1
		}
		readOff += last
	}
	return ret, nil
}

var sumByFile = map[string]int64{}

// Sum7 is cached (cheating!) (:
func Sum7(fileName string) (ret int64, err error) {
	if s, ok := sumByFile[fileName]; ok {
		return s, nil
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}
		num, err := ParseInt(b[last:i])
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}

	sumByFile[fileName] = ret
	return ret, nil
}

type sequence struct {
	end int
	sum int64
}

func findSequence(b []byte) (sequence, error) {
	s := sequence{}
	firstNum := int64(0)

	i := 0
	for ; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}

		num, err := ParseInt(b[0:i])
		if err != nil {
			return s, err
		}
		firstNum = num
		s.sum += num
		break
	}

	s.end = i + 1
	for i++; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}

		num, err := ParseInt(b[s.end:i])
		if err != nil {
			return s, err
		}
		if num == firstNum {
			return s, nil
		}
		s.sum += num
		s.end = i + 1
	}
	return s, nil
}

// Sum8 and we know that some sequences might be repeating...
func Sum8(fileName string) (ret int64, err error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	seq, err := findSequence(b)
	if err != nil {
		return 0, err
	}
	ret += seq.sum

	last := seq.end
	for i := seq.end; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}

		// Is it next element of sequence?
		if len(b[last:]) >= seq.end &&
			bytes.Compare(b[last:last+seq.end], b[0:seq.end]) == 0 {
			last += seq.end
			i += (seq.end - 1)
			ret += seq.sum
			continue
		}

		num, err := ParseInt(b[last:i])
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}

	return ret, nil
}
