package main

import (
	"fmt"
	service2 "github.com/hongkailiu/test-go/pkg/experimental/test/testify/service"
)

func main() {
	d := service2.NewDB()

	g := service2.NewGreeter(d, "en")
	fmt.Println(g.Greet())             // Message is: hello
	fmt.Println(g.GreetInDefaultMsg()) // Message is: default message

	g = service2.NewGreeter(d, "es")
	fmt.Println(g.Greet()) // Message is: holla

	g = service2.NewGreeter(d, "random")
	fmt.Println(g.Greet()) // Message is: bzzzz
}
