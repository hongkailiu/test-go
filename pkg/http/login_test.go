package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLocalID(t *testing.T) {
	u := user{}
	assert.Empty(t, u.localID)
	u.setLocalID()
	assert.NotEmpty(t, u.localID)
}