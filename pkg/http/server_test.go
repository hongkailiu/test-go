package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hongkailiu/test-go/pkg/http/db"
	cmdconfig "github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/stretchr/testify/assert"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func TestRoute1(t *testing.T) {

	appConfig = loadConfig()

	hc := cmdconfig.HttpConfig{}

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
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
