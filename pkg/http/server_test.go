package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/diff"

	"github.com/hongkailiu/test-go/pkg/http/info"
	"github.com/hongkailiu/test-go/pkg/http/model"
	cmdconfig "github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type MyMockedDBService struct {
	mock.Mock
}

func (m *MyMockedDBService) GetCities(limit, offset int) (*[]model.City, error) {
	args := m.Called(limit, offset)
	return args.Get(0).(*[]model.City), args.Error(1)

}

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

	mock := new(MyMockedDBService)
	cities := []model.City{{Name: "test-city", Model: gorm.Model{ID: uint(23)}}}
	mock.On("GetCities", 10, 0).Return(&cities, nil)

	router := setupRouter(&hc, githubLogin, googleLogin, mock)

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

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/v1/cities", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+testToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//var cities *[]model.City
	err = json.Unmarshal([]byte(w.Body.Bytes()), &cities)
	assert.Nil(t, err)
	expected := []model.City{{Name: "test-city", Model: gorm.Model{ID: uint(23)}}}
	if !reflect.DeepEqual(expected, cities) {
		t.Errorf("Unexpected mis-match: %s", diff.ObjectReflectDiff(expected, cities))
	}
}

func TestSetupOAuthConfig1(t *testing.T) {
	appConfig = loadConfig()

	oauthConfGitHub, err := setupOAuthConfig("github")
	assert.Nil(t, err)
	assert.NotNil(t, oauthConfGitHub)

	oauthConfGoogle, err := setupOAuthConfig("google")
	assert.Nil(t, err)
	assert.NotNil(t, oauthConfGoogle)

	oauthConfAWS, err := setupOAuthConfig("aws")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("do not support oauth for provider %s", "aws"), err.Error())
	assert.Nil(t, oauthConfAWS)
}
