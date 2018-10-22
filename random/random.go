package random

import (
	"math/rand"
	"time"
)

func GetRandom(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
