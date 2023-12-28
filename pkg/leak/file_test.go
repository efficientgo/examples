// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package leak

import (
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
)

func TestDoWithFile(t *testing.T) {
	testutil.Ok(t, doWithFile_Wrong("/dev/null"))
	testutil.Ok(t, doWithFile_CaptureCloseErr("/dev/null"))
	doWithFile_LogCloseErr(log.NewNopLogger(), "/dev/null")
}

func TestOpenMultiple(t *testing.T) {
	files, err := openMultiple_Wrong("/dev/null", "/dev/null", "/dev/null")
	testutil.Ok(t, err)
	testutil.Ok(t, closeAll(files))

	files, err = openMultiple_Correct("/dev/null", "/dev/null", "/dev/null")
	testutil.Ok(t, err)
	testutil.Ok(t, closeAll(files))
}
