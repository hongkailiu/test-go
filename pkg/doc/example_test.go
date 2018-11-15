package doc_test

import (
	"fmt"
	"github.com/hongkailiu/test-go/pkg/doc"
)

func Example() {
	fmt.Println("It is pkg level")
	// Output: It is pkg level
}

func ExampleMYNewNiceHello_canada() {
	doc.MYNewNiceHello()
	fmt.Println("new hello canada")
}

func ExampleStart() {
	fmt.Println("it started a function!")
	// Output: it started a function!
}

func ExampleStart_warm() {
	fmt.Println("it started warm!")
	// Output: it started warm!
}

func ExampleStart_cold() {
	fmt.Println("it started cold!")
	// Output: it started cold!
}
