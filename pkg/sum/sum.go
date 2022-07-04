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

// ParseInt is 3-4x times faster than strconv.ParseInt or Atoi.
func ParseInt(input []byte) (n int64, _ error) {
	factor := int64(1)
	k := 0
	// TODO(bwplotka): Optimize if only positive integers are accepted.
	if input[0] == '-' {
		factor *= -1
		k++
	}

	for i := len(input) - 1; i >= k; i-- {
		if input[i] < '0' || input[i] > '9' {
			return 0, errors.Errorf("not a valid integerer: %v", input)
		}

		n += factor * int64(input[i]-'0')
		factor *= 10
	}
	return n, nil
}

func Sum2(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var last int
	for begin := 0; begin < len(b); begin++ {
		if b[begin] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(string(b[last:begin]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = begin + 1
	}
	return ret, nil
}

func zeroCopyToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

func Sum3(fileName string) (ret int64, err error) {
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

func Sum4(fileName string) (ret int64, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer errcapture.Do(&err, f.Close, "close file")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		num, err := strconv.ParseInt(zeroCopyToString(scanner.Bytes()), 10, 64)
		// Or just use our custom int parser function we used in.
		// num, err := ParseInt(scanner.Bytes())
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil
}
