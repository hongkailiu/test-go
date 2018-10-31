package http

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	googleauth "google.golang.org/api/oauth2/v2"
	"net/http"
	"strconv"
)

var (
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = uuid.NewV4().String()

	options = sessions.Options{
		//Path:     "/",
		//Domain:   "hongkailiu.tk",
		MaxAge: 300,
		//Secure:   false,
		//HttpOnly: false,
	}
)

func saveDataInSession(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Options(options)
	session.Set("username", username)
	session.Save()
}

type user struct {
	name  string
	id    string
	email string
}

// UserProvider provides user
type UserProvider interface {
	getUser(client *http.Client) (*user, error)
}

type login struct {
	config       *oauth2.Config
	userProvider UserProvider
}

func (l login) getLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := l.config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		log.WithFields(log.Fields{"url": url}).Debug("redirect login url")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (l login) getCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		code := r.FormValue("code")
		token, err := l.config.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		u, err := l.userProvider.getUser(l.config.Client(oauth2.NoContext, token))
		if err != nil {
			fmt.Printf("l.userProvider.getUser() failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}
		//fmt.Printf("Logged in as google user: %s\n", (*u).name)
		log.Debugf("Logged in as google user: %v", u)
		log.WithFields(log.Fields{"u.name": u.name}).Debug("saving u.name")
		saveDataInSession(c, u.name)
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	}
}

type googleUserProvider struct {
}

func (up googleUserProvider) getUser(client *http.Client) (*user, error) {
	googleService, err := googleauth.New(client)
	if err != nil {
		return nil, err
	}
	userinfoplus, err := googleService.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}
	return &user{name: (*userinfoplus).Name, id: (*userinfoplus).Id, email: (*userinfoplus).Email}, nil
}

type gitHubUserProvider struct {
}

func (up gitHubUserProvider) getUser(client *http.Client) (*user, error) {
	gitHubClient := github.NewClient(client)
	u, _, err := gitHubClient.Users.Get(context.Background(), "")
	log.Debugf("get github user: %v", u)
	if err != nil {
		return nil, err
	}
	//https://stackoverflow.com/questions/35373995/github-user-email-is-null-despite-useremail-scope
	result := &user{name: *(u.Name), id: strconv.FormatInt(*(u.ID), 10)}
	if u.Email != nil {
		result.email = *(u.Email)
		return result, nil
	}

	userEmails, _, err := gitHubClient.Users.ListEmails(context.Background(), &github.ListOptions{Page: 1, PerPage: 30})
	if err != nil {
		return nil, err
	}

	for _, email := range userEmails {
		if email != nil && email.Primary != nil && *email.Primary {
			result.email = *email.Email
			return result, nil
		}

	}

	result.email = ""
	return result, nil
}
