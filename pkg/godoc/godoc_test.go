package godoc_test

import (
	"fmt"
	"testing"

	. "github.com/hongkailiu/test-go/pkg/godoc"
)

func TestNewHello(t *testing.T) {
	err := NewHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}

func TestNewNiceHello(t *testing.T) {
	err := NewNiceHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}

func ExampleNewNiceHello_Cool() {
	//More info https://github.com/fluhus/godoc-tricks/blob/master/doc.go
	NewNiceHello()
	fmt.Println("It is cool")
}

func ExampleNewNiceHello_Warm() {
	NewNiceHello()
	fmt.Println("It is warm")
}
