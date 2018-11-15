package godoc_test

import (
	"fmt"

	"github.com/hongkailiu/test-go/pkg/godoc"
)

func ExampleNewNiceHello_Cool() {
	//More info https://github.com/fluhus/godoc-tricks/blob/master/doc.go
	godoc.NewNiceHello()
	fmt.Println("It is cool")
}

func ExampleNewNiceHello_Warm() {
	godoc.NewNiceHello()
	fmt.Println("It is warm")
}
