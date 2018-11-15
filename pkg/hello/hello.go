package main

import (
	"fmt"
	"github.com/hongkailiu/test-go/pkg/stringutil"
)

func main() {
	//fmt.Printf("hello, world\n")
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
}

//Deprecated: NewHello is defined here to show
//the usage for `Deprecated` in goDoc
func NewHello() error {
	return nil
}