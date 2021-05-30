package main

func use(_ int) {}

func main() {
	var a int // error: a declared but not used <1>

	b := 1 // error: b declared but not used <1>

	var c int
	d := c // error: d declared but not used <1>

	e := 1
	use(e)

	f := 1
	_ = f
}
