package http

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	defaultSessionKeyForTest = []byte{116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116}
)

func TestLoadConfig1(t *testing.T) {
	log = logrus.New()
	c := loadConfig()
	assert.Equal(t, defaultSessionKeyForTest, c.sessionKey)
}

func TestLoadConfig2(t *testing.T) {
	log = logrus.New()
	os.Setenv("session_key", "[")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKeyForTest, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}

func TestLoadConfig3(t *testing.T) {
	log = logrus.New()
	os.Setenv("session_key", "[116 117 118]")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKeyForTest, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}

func TestLoadConfig4(t *testing.T) {
	log = logrus.New()
	os.Setenv("session_key", "[116 117 118 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116]")
	c := loadConfig()
	assert.Equal(t, []byte{116, 117, 118, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116, 116}, c.sessionKey)
}

func TestLoadConfig5(t *testing.T) {
	log = logrus.New()
	os.Setenv("session_key", "[116 aa 118 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116 116]")
	c := loadConfig()
	assert.NotEqual(t, defaultSessionKeyForTest, c.sessionKey)
	assert.Equal(t, 32, len(c.sessionKey))
}

func TestLoadDBConfig1(t *testing.T) {
	log = logrus.New()
	c := loadDBConfig()
	assert.Equal(t, "host=localhost port=5432 user=redhat dbname=ttt password=redhat sslmode=disable", c.getDBString())
}
