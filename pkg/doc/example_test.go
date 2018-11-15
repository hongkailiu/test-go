package doc_test

import (
	"fmt"
)


func Example() {
	fmt.Println("It is pkg level")
}

func ExampleStart() {
	fmt.Println("it started!")
	// Output: it started!
}

func ExampleStart_Warm() {
	fmt.Println("it started warm!")
	// Output: it started warm!
}

func ExampleStart_Cold() {
	fmt.Println("it started cold!")
	// Output: it started cold!
}
