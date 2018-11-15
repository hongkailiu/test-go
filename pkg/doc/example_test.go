package doc_test

import (
	"fmt"
)


func Example() {
	fmt.Println("It is pkg level")
}


func ExampleStart_warm() {
	fmt.Println("it started warm!")
	// Output: it started warm!
}

func ExampleStart_cold() {
	fmt.Println("it started cold!")
	// Output: it started cold!
}
