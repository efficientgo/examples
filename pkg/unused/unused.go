// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

// Example of compilation errors from unused variables.
// Read more in "Efficient Go"; Example 2-7.

func use(_ int) {}

func main() {
	// var a int // error: a declared but not used

	// b := 1 // error: b declared but not used

	// var c int
	// d := c // error: d declared but not used

	e := 1
	use(e)

	f := 1
	_ = f
}
