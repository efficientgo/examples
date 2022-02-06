package main

import (
	"fmt"

	"github.com/efficientgo/examples/pkg/compiler/packageA"
	"github.com/efficientgo/examples/pkg/compiler/packageB"
)

func main() {
	const a = 6

	fmt.Println(packageA.A(a) + packageB.B(a))
}
