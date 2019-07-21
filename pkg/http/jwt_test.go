package http

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	testLocalID = "localID"
	testToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjM2ODc1NzksImp0aSI6ImxvY2FsSUQifQ.HZCSw27dRZMQjIBvEgbpZO5rUVrO0LrCwWH9MFG4J-M"
)

func TestJWT1(t *testing.T) {
	token, err := getToken(testLocalID, defaultSessionKeyForTest)
	assert.Nil(t, err)
	localID, err := getLocalIDFromToken(token, defaultSessionKeyForTest)
	assert.Nil(t, err)
	assert.Equal(t, testLocalID, localID)
}

func TestJWT2(t *testing.T) {
	token, err := getExpiredToken(testLocalID, defaultSessionKeyForTest)
	assert.Nil(t, err)
	localID, err := getLocalIDFromToken(token, defaultSessionKeyForTest)
	assert.Empty(t, localID)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "token is expired"))
}

func TestJWT3(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)

	token, err := getNonMetheodToken(testLocalID, key)
	assert.Nil(t, err)
	localID, err := getLocalIDFromToken(token, defaultSessionKeyForTest)
	assert.Empty(t, localID)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "unexpected signing alg"))
}

func getExpiredToken(localID string, key interface{}) (string, error) {
	claims := &jwt.StandardClaims{
		//https://godoc.org/github.com/dgrijalva/jwt-go#pkg-examples
		Id:        localID,
		ExpiresAt: time.Now().Add(-time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func getNonMetheodToken(localID string, key interface{}) (string, error) {
	claims := &jwt.StandardClaims{
		//https://godoc.org/github.com/dgrijalva/jwt-go#pkg-examples
		Id:        localID,
		ExpiresAt: time.Now().Add(time.Minute * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(key)
}
