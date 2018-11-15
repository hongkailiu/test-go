package doc_test

import (
	"fmt"

	"github.com/hongkailiu/test-go/pkg/doc"
)

func ExampleMYNewNiceHello_USA() {
	//More info https://github.com/fluhus/godoc-tricks/blob/master/doc.go
	doc.MYNewNiceHello()
	fmt.Println("It is USA")
}

func ExampleMYNewNiceHello_Canada() {
	doc.MYNewNiceHello()
	fmt.Println("It is Canada")
}

func Example() {
	fmt.Println("It is pkg level")
}

func ExampleStart() {
	fmt.Println("it started!")
	// Output: dummy
}

func ExampleStart_Warm() {
	fmt.Println("it started warm!")
	// Output: dummy
}

func ExampleStart_Cold() {
	fmt.Println("it started cold!")
	// Output: dummy
}
