package main

import "github.com/pkg/errors"

func shouldFail() bool { return false }

func noErrCanHappen() int {
	// ...
	return 204
}

func doOrErr() error {
	// ...
	if shouldFail() {
		return errors.New("ups, XYZ failed")
	}
	return nil
}

func intOrErr() (int, error) {
	// ...
	if shouldFail() {
		return 0, errors.New("ups, XYZ2 failed")
	}
	return noErrCanHappen(), nil
}

func main() {
	ret := noErrCanHappen()
	if err := nestedDoOrErr(); err != nil {
		// handle error
	}
	ret2, err := intOrErr()
	if err != nil {
		// handle error
	}
	// ...
	_, _ = ret, ret2 // Just so we can compile the code.
}

func nestedDoOrErr() error {
	// ...
	if err := doOrErr(); err != nil {
		return errors.Wrap(err, "do")
	}
	return nil
}