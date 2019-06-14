package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TTT struct { i int; f float64; next *TTT }

func TestType1(t *testing.T) {
	var objectUnderTest TTT
	expected := TTT{i:0, f:0.0, next:nil}
	assert.Equal(t, expected, objectUnderTest, "they should be equal")
}
