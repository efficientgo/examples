package sum

import (
	"bufio"
	"bytes"
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
	for curr := 0; curr < len(b); curr++ {
		if b[curr] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(string(b[last:curr]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = curr + 1
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
	for curr := 0; curr < len(b); curr++ {
		if b[curr] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(zeroCopyToString(b[last:curr]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = curr + 1
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
	for curr := 0; curr < len(b); curr++ {
		if b[curr] != '\n' {
			continue
		}
		num, err := ParseInt(b[last:curr])
		if err != nil {
			return 0, err
		}

		ret += num
		last = curr + 1
	}
	return ret, nil
}

// Slower than Sum4.
func Sum4_scanner(fileName string) (ret int64, err error) {
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
