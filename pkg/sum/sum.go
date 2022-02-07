package sum

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

func Sum(fileName string) (ret int64, _ error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	for _, line := range bytes.Split(b, []byte("\n")) {
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
		num, err := ParseInt(b[last:begin])
		if err != nil {
			// TODO(bwplotka): Return err using other channel.
			continue
		}
		ret += num
		last = begin + 1
	}
	return ret, nil
}
