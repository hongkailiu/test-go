package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hongkailiu/test-go/pkg/http/db"
	"github.com/hongkailiu/test-go/pkg/http/info"
	cmdconfig "github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/stretchr/testify/assert"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func TestRoute1(t *testing.T) {

	appConfig = loadConfig()

	hc := cmdconfig.HttpConfig{Version: "test-version"}

	oauthConfGitHub := &oauth2.Config{
		ClientID:     "appConfig.ghClientID",
		ClientSecret: "appConfig.ghClientSecret",
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     githuboauth.Endpoint,
	}

	//https://console.developers.google.com/apis/dashboard
	oauthConfGoogle := &oauth2.Config{
		ClientID:     "appConfig.ggClientID",
		ClientSecret: "appConfig.ggClientSecret",
		RedirectURL:  "appConfig.ggRedirectURL",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	githubLogin := login{oauthConfGitHub, gitHubUserProvider{}}
	googleLogin := login{oauthConfGoogle, googleUserProvider{}}

	dbService := db.Service{}

	router := setupRouter(&hc, githubLogin, googleLogin, &dbService)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	i := info.Info{}
	err = json.Unmarshal([]byte(w.Body.Bytes()), &i)
	assert.Nil(t, err)
	assert.Equal(t, hc.Version, i.Version)
	assert.NotEmpty(t, i.Ips)
	assert.True(t, time.Now().Sub(i.Now) > 0)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/v1/users", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/token", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
