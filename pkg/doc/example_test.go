package doc_test

import (
	"fmt"

	"github.com/hongkailiu/test-go/pkg/doc"
)

func ExampleNewNiceHello_Cool() {
	//More info https://github.com/fluhus/godoc-tricks/blob/master/doc.go
	doc.NewNiceHello()
	fmt.Println("It is cool")
}

func ExampleNewNiceHello_Warm() {
	doc.NewNiceHello()
	fmt.Println("It is warm")
}

func Example() {
	fmt.Println("It is pkg level")
}

func ExampleStart() {
	fmt.Println("it started!")
	// Output: dummy
}
