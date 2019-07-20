package http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultSessionKey = []byte{116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116}
)

func TestLoadConfig1(t *testing.T) {
	c := loadConfig()
	assert.Equal(t, defaultSessionKey, c.sessionKey)
}

func TestLoadConfig2(t *testing.T) {
	os.Setenv("session_key", "[")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKey, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}

func TestLoadConfig3(t *testing.T) {
	os.Setenv("session_key", "[116 117 118]")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKey, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}

func TestLoadConfig4(t *testing.T) {
	os.Setenv("session_key", "[116 117 118 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116]")
	c := loadConfig()
	assert.Equal(t, []byte{116, 117, 118, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116}, c.sessionKey)
}

func TestLoadConfig5(t *testing.T) {
	os.Setenv("session_key", "[116 aa 118 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116]")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKey, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}
