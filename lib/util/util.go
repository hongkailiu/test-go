package util

import (
	"math/rand"
	"time"
)

// GetRandomInt return a random integer
func GetRandomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}
