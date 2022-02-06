package sum

import (
	"bytes"
	"io/ioutil"
	"strconv"
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
