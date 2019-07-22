package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetLocalID(t *testing.T) {
	u := user{}
	assert.Empty(t, u.localID)
	u.setLocalID()
	assert.NotEmpty(t, u.localID)
}