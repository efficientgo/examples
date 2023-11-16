package sum

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"unsafe"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
)

// Sum is a naive implementation and algorithm for summing integers from file.
// Read more in "Efficient Go"; Example 4-1.
func Sum(fileName string) (ret int64, _ error) {
	b, err := os.ReadFile(fileName)
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
// Read more in "Efficient Go"; Example 10-3.
func Sum2(fileName string) (ret int64, _ error) {
	b, err := os.ReadFile(fileName)
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

func zeroCopyToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

// Sum3 is a sum with optimized the second latency + CPU bottleneck: string conversion.
// On CPU profile we see byte to string conversion not only allocate memory, but also takes precious time.
// Let's perform zeroCopy conversion.
// 2x less latency memory than Sum2.
// Read more in "Efficient Go"; Example 10-4.
func Sum3(fileName string) (ret int64, _ error) {
	b, err := os.ReadFile(fileName)
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

// ParseInt is 3-4x times faster than strconv.ParseInt or Atoi.
func ParseInt(input []byte) (n int64, _ error) {
	factor := int64(1)
	k := 0

	// TODO(bwplotka): Optimize if only positive integers are accepted (only 2.6% overhead in my tests though).
	if input[0] == '-' {
		factor *= -1
		k++
	}

	for i := len(input) - 1; i >= k; i-- {
		if input[i] < '0' || input[i] > '9' {
			return 0, errors.Newf("not a valid integer: %v", input)
		}

		n += factor * int64(input[i]-'0')
		factor *= 10
	}
	return n, nil
}

// Sum4 is a sum with optimized the second latency + CPU bottleneck: ParseInt and string conversion.
// On CPU profile we see that ParseInt does a lot of checks that we might not need. We write our own parsing
// straight from byte to avoid conversion CPU time.
// 2x less latency, same mem as Sum3.
// Read more in "Efficient Go"; Example 10-5.
func Sum4(fileName string) (ret int64, err error) {
	b, err := os.ReadFile(fileName)
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

func Sum4_atoi(fileName string) (ret int64, err error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ {
		if b[i] != '\n' {
			continue
		}
		num, err := strconv.Atoi(zeroCopyToString(b[last:i]))
		if err != nil {
			return 0, err
		}

		ret += int64(num)
		last = i + 1
	}
	return ret, nil
}

// Sum5 is like Sum4, but noticing that it takes time to even allocate 21 MB on heap (and read file to it).
// Let's try to use scanner instead.
// Slower than Sum4 and Sum6 because scanner is not optimized for this...? Scanner takes 73% of CPU time.
// Read more in "Efficient Go"; Example 10-7.
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
	return ret, scanner.Err()
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
	return ret, scanner.Err()
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
// Assuming no integer is larger than 8 000 digits.
// Read more in "Efficient Go"; Example 10-8.
func Sum6(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	buf := make([]byte, 8*1024)
	return Sum6Reader(f, buf)
}

func Sum6Reader(r io.Reader, buf []byte) (ret int64, err error) { // Just inlining this function saves 7% on latency
	var offset, n int
	for err != io.EOF {
		n, err = r.Read(buf[offset:])
		if err != nil && err != io.EOF {
			return 0, err
		}
		n += offset

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

		offset = n - last
		if offset > 0 {
			_ = copy(buf, buf[last:n])
		}
	}
	return ret, nil
}

var sumByFile = map[string]int64{}

// Sum7 is cached (cheating!) (:
// Read more in "Efficient Go"; Example 10-15.
func Sum7(fileName string) (int64, error) {
	if s, ok := sumByFile[fileName]; ok {
		return s, nil
	}

	ret, err := Sum(fileName)
	if err != nil {
		return 0, err
	}

	sumByFile[fileName] = ret
	return ret, nil
}
