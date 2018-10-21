package main

import (
	"github.com/hongkailiu/test-go/http"
)

func main() {
	http.Server{Port:8080}.Run()
}