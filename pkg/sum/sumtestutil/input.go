package sumtestutil

import (
	"io"

	"github.com/pkg/errors"
)

func CreateTestInputWithExpectedResult(w io.Writer, numLen int) (sum int64, err error) {
	const testSumOfTen = int64(31108)
	var tenSet = []byte(`123
43
632
22
2
122
26660
91
2
3411
`)

	if numLen%10 != 0 {
		return 0, errors.Errorf("number of input should be division by 10, got %v", numLen)
	}

	for i := 0; i < numLen/10; i++ {
		if _, err := w.Write(tenSet); err != nil {
			return 0, err
		}
	}

	return testSumOfTen * (int64(numLen) / 10), nil
}
