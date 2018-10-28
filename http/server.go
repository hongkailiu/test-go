package http

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/hongkailiu/test-go/random"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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
)

// PrometheusLogger intercepts all http requests and logging the path
func PrometheusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"c.Request.URL.Path": c.Request.URL.Path,
		}).Debug("prometheus logger detected path visited")
		httpRequestsTotal.With(prometheus.Labels{"path": c.Request.URL.Path}).Inc()
	}
}

// Run starts the http server
func Run() {

	log.WithFields(log.Fields{"oauthConf.ClientID": oauthConf.ClientID, "oauthConf.ClientSecret": oauthConf.ClientSecret}).Debug("oauthConf")

	prometheusRegister()

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	gin.Logger()
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(PrometheusLogger())

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", func(c *gin.Context) {
		infoP := getInfo()
		c.JSON(http.StatusOK, *infoP)
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	staticDir := filepath.Join(dir, "static")
	log.WithFields(log.Fields{"staticDir": staticDir}).Debug("http staticDir dir")
	r.StaticFS("/console", http.Dir(staticDir))

	//https://blog.kowalczyk.info/article/f/accessing-github-api-from-go.html
	r.GET("/login", func(c *gin.Context) {
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		log.WithFields(log.Fields{"url": url}).Debug("redirect login url")
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/github_oauth_cb", func(c *gin.Context) {
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
		fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
		session := sessions.Default(c)
		session.Set("username", *user.Login)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	})

	r.GET("/whoami", func(c *gin.Context) {
		session := sessions.Default(c)
		var username string
		v := session.Get("username")
		if v == nil {
			username = ""
		} else {
			username = v.(string)
		}
		c.JSON(200, gin.H{"username": username})
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("username")
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	})

	go func() {
		for {
			n := random.GetRandom(1000)
			log.WithFields(log.Fields{"n": n}).Debug("generated random number")
			randomNumber.With(prometheus.Labels{"key": "value"}).Set(float64(n))
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			n := random.GetRandom(50)
			log.WithFields(log.Fields{
				"n": n,
			}).Debug("generated fake storageOperationMetric")
			storageOperationMetric.With(prometheus.Labels{"volume_plugin": "hongkailiu.tk/aws-ebs", "operation_name": "volume_provision"}).Observe(float64(n))
			time.Sleep(100 * time.Second)
		}
	}()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}
