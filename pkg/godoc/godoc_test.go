package godoc_test

import (
	"testing"

	. "github.com/hongkailiu/test-go/pkg/godoc"
)

func TestNewHello(t *testing.T) {
	err := NewHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}

// You can place usage examples in your godoc.
//
// Examples should be placed in a file with a _test suffix. For example, the
// examples in this guide are in a file called godoc_test.go .
//
// The example functions should be called Example() for
// package examples, ExampleTypename() for a specific type or
// ExampleFuncname() for a specific function. For multiple examples
// for the same entity (like same function), you can add a suffix like
// ExampleFoo_suffix1, ExampleFoo_suffix2.
//
// You can document an example's output, by adding an output comment at its end.
// The output comment must begin with "Output:", as shown below:
//  func ExampleNewNiceHello() {
//      fmt.Println("Hello")
//      // Output: Hello
//  }
//
// Notice that the tricks brought here (titles, code blocks, links etc.) don't work
// in example documentation.
func TestNewNiceHello(t *testing.T) {
	err := NewNiceHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}
