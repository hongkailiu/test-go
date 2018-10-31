package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	googleauth "google.golang.org/api/oauth2/v2"
)

var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
	// Set ClientId and ClientSecret to
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("gh_client_id"),
		ClientSecret: os.Getenv("gh_client_secret"),
		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"user:email"},
		Endpoint: githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = uuid.NewV4().String()

	//https://console.developers.google.com/apis/dashboard
	oauthConfGoogle = &oauth2.Config{
		ClientID:     os.Getenv("gg_client_id"),
		ClientSecret: os.Getenv("gg_client_secret"),
		RedirectURL:  "http://127.0.0.1:8080/google_oauth_cb",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	options = sessions.Options{
		//Path:     "/",
		//Domain:   "hongkailiu.tk",
		MaxAge: 300,
		//Secure:   false,
		//HttpOnly: false,
	}

	githubLoginHandler = func(c *gin.Context) {
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		log.WithFields(log.Fields{"url": url}).Debug("redirect login url")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}

	githubCallbackHandler = func(c *gin.Context) {
		r := c.Request
		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		code := r.FormValue("code")
		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		client := github.NewClient(oauthClient)
		user, _, err := client.Users.Get(context.Background(), "")
		if err != nil {
			fmt.Printf("client.Users.Get() faled with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}
		fmt.Printf("Logged in as GitHub user: %s\n", *user.Name)
		saveDataInSession(c, *user.Name)
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	}

	googleLoginHandler = func(c *gin.Context) {
		url := oauthConfGoogle.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		log.WithFields(log.Fields{"url": url}).Debug("redirect login url")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}

	googleCallbackHandler = func(c *gin.Context) {
		r := c.Request
		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		code := r.FormValue("code")
		token, err := oauthConfGoogle.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}

		oauthClient := oauthConfGoogle.Client(oauth2.NoContext, token)
		googleService, err := googleauth.New(oauthClient)
		if err != nil {
			fmt.Printf("googleauth.New(oauthClient) faled with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}
		userinfoplus, err := googleService.Userinfo.Get().Do()
		if err != nil {
			fmt.Printf("googleService.Userinfo.Get().Do() faled with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/console")
			return
		}
		fmt.Printf("Logged in as google user: %s\n", (*userinfoplus).Name)
		saveDataInSession(c, (*userinfoplus).Name)
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	}
)

func saveDataInSession(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Options(options)
	session.Set("username", username)
	session.Save()
}
