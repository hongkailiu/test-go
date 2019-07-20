package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testLocalID = "localID"
	testToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjM2ODc1NzksImp0aSI6ImxvY2FsSUQifQ.HZCSw27dRZMQjIBvEgbpZO5rUVrO0LrCwWH9MFG4J-M"
)

func TestJWT1(t *testing.T) {
	token, err := getToken(testLocalID, defaultSessionKey)
	assert.Nil(t, err)
	localID, err := getLocalIDFromToken(token, defaultSessionKey)
	assert.Nil(t, err)
	assert.Equal(t, testLocalID, localID)
}
