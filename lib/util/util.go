package util

import (
	"math/rand"
	"time"
)

func GetRandomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}