package http

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

type config struct {
	sessionKey     []byte
	ghClientID     string
	ghClientSecret string
	ggClientID     string
	ggClientSecret string
	ggRedirectURL  string
}

func loadConfig() *config {
	config := &config{}
	config.sessionKey = securecookie.GenerateRandomKey(32)
	newSessionKey, err := getSessionKey()
	if err != nil {
		log.Warnf("error found when getSessionKey(): %s", err.Error())
	} else {
		config.sessionKey = newSessionKey
	}
	config.ghClientID = os.Getenv("gh_client_id")
	config.ggClientSecret = os.Getenv("gh_client_secret")
	config.ggClientID = os.Getenv("gg_client_id")
	config.ggClientSecret = os.Getenv("gg_client_secret")
	config.ggRedirectURL = os.Getenv("gg_redirect_url")

	if config.ggRedirectURL == "" {
		config.ggRedirectURL = "http://127.0.0.1:8080/google_oauth_cb"
	}
	return config
}

func getSessionKey() ([]byte, error) {
	key := os.Getenv("session_key")
	if key == "" {
		return nil, fmt.Errorf("env. var. session_key not found")
	}
	trimKey := strings.TrimSuffix(strings.TrimPrefix(key, "["), "]")
	bytes := strings.Split(trimKey, " ")
	if len(bytes) != 32 {
		return nil, fmt.Errorf("key length is %d (only 32 allowed)", len(bytes))
	}

	var result []byte
	for _, b := range bytes {
		i, err := strconv.ParseUint(b, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("get error when strconv.ParseUint(b, 10, 8) for b=%s", b)
		}
		result = append(result, byte(i))
	}
	log.WithFields(log.Fields{"result": result}).Warnf("got secret from env. var. session_key")
	return result, nil
}
