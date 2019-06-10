package json

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypeDuration(t *testing.T) {
	objectUnderTest := 3*time.Second + 5*time.Minute
	actual := fmt.Sprintf("%s", objectUnderTest)
	assert.Equal(t, "5m3s", actual, "they should be equal")
}

func TestTypeDuration1(t *testing.T) {
	td := TestDuration{
		TestDuration: 23*time.Second + 6*time.Minute,
	}

	bolB, err := json.Marshal(td)
	assert.Nil(t, err)

	actual := string(bolB)

	//fmt.Println("=" + actual + "=")

	expect := `{"marker":"","name":"","type":"","startTime":"0001-01-01T00:00:00Z","testDuration":383000000000,"Steps":null}`

	assert.Equal(t, expect, actual, "they should be equal")
}
