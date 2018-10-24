package util

import (
	"math/rand"
	"time"
)

// GetRandomInt return a random integer
func GetRandomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}


// PanicIfError panics if error is not nil
func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}