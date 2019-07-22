package http

import (
	"encoding/json"
	"fmt"
	"github.com/hongkailiu/test-go/pkg/swagger/swagger/models"
	"net/http"
	"net/http/httptest"
	"os"
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

var (
	hc                       cmdconfig.HttpConfig
	githubLogin, googleLogin login
)

type MyMockedDBService struct {
	mock.Mock
}

func (m *MyMockedDBService) GetCities(limit, offset int) (*[]model.City, error) {
	args := m.Called(limit, offset)
	return args.Get(0).(*[]model.City), args.Error(1)

}

func beforeEach() {
	fmt.Println("============================beforeEach======================")
	appConfig = loadConfig()

	hc = cmdconfig.HttpConfig{Version: "test-version"}

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

	githubLogin = login{oauthConfGitHub, gitHubUserProvider{}}
	googleLogin = login{oauthConfGoogle, googleUserProvider{}}
}

func TestRoute1(t *testing.T) {
	beforeEach()
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
	err = json.Unmarshal(w.Body.Bytes(), &cities)
	assert.Nil(t, err)
	expected := []model.City{{Name: "test-city", Model: gorm.Model{ID: uint(23)}}}
	if !reflect.DeepEqual(expected, cities) {
		t.Errorf("Unexpected mis-match: %s", diff.ObjectReflectDiff(expected, cities))
	}
}

func TestRoute2(t *testing.T) {
	beforeEach()
	mock := new(MyMockedDBService)
	cities := []model.City{{Name: "test-city", Model: gorm.Model{ID: uint(23)}}}
	errorMsg := "test-error"
	mock.On("GetCities", 10, 0).Return(&cities, fmt.Errorf(errorMsg))

	router := setupRouter(&hc, githubLogin, googleLogin, mock)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/whoami", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+testToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Header().Get("username"))

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/token", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+testToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "", w.Header().Get("token"))

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/v1/cities", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+testToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var mError models.Error
	err = json.Unmarshal(w.Body.Bytes(), &mError)
	assert.Nil(t, err)
	expected := models.Error{Code: int64(http.StatusInternalServerError), Message: &errorMsg}
	if !reflect.DeepEqual(expected, mError) {
		t.Errorf("Unexpected mis-match: %s", diff.ObjectReflectDiff(expected, mError))
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

func TestBeforeStartServer(t *testing.T) {
	hc = cmdconfig.HttpConfig{Version: "test-version"}
	os.Setenv("unit_testing", "true")
	Run(&hc)
}

func TestGetSecret(t *testing.T) {
	assert.Len(t, GetSecret(16), 16)
}
