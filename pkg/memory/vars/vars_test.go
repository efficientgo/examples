package vars

import (
	"fmt"
	"testing"
	"unsafe"
)

func test() {
	type A struct { //<1>
		Label1 int
		Label2 string
		Label3 []int
	}

	var label1 A                  // <2>
	label1.Label2 = "some string" // <3>

	var label2, label3, label4 *A    // <2>
	label2 = &A{Label3: []int{1, 2}} // <4>
	label3 = label2                  // <5>
	label4 = &label1                 // <5>

	fmt.Println(label1, label2, label3, label4, unsafe.Sizeof(A{}))
	fmt.Printf("%T", label2)
}

func TestTest(t *testing.T) {
	test()
}
