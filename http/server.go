package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/securecookie"
	"github.com/hongkailiu/test-go/random"
	"github.com/hongkailiu/test-go/swagger/swagger/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
	// Set ClientId and ClientSecret to
	oauthConfGitHub = &oauth2.Config{
		ClientID:     os.Getenv("gh_client_id"),
		ClientSecret: os.Getenv("gh_client_secret"),
		// select level of access you want from https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"read:user", "user:email"},
		Endpoint: githuboauth.Endpoint,
	}

	//https://console.developers.google.com/apis/dashboard
	oauthConfGoogle = &oauth2.Config{
		ClientID:     os.Getenv("gg_client_id"),
		ClientSecret: os.Getenv("gg_client_secret"),
		RedirectURL:  "http://127.0.0.1:8080/google_oauth_cb",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	githubLogin = login{oauthConfGitHub, gitHubUserProvider{}}
	googleLogin = login{oauthConfGoogle, googleUserProvider{}}
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

	log.WithFields(log.Fields{"oauthConfGitHub.ClientID": oauthConfGitHub.ClientID, "oauthConfGitHub.ClientSecret": oauthConfGitHub.ClientSecret}).Debug("oauthConf")
	log.WithFields(log.Fields{"oauthConfGoogle.ClientID": oauthConfGoogle.ClientID, "oauthConfGoogle.ClientSecret": oauthConfGoogle.ClientSecret}).Debug("oauthConf")

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

	// sessions and cookies are from github.com/gin-contrib
	// which uses implementation from github.com/gorilla/
	sessionKey := securecookie.GenerateRandomKey(32)
	log.WithFields(log.Fields{"sessionKey": sessionKey}).Info("generated session key")
	store := cookie.NewStore(sessionKey)
	r.Use(sessions.Sessions("my_session", store))

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
	r.GET("/login_github", githubLogin.getLoginHandler())
	r.GET("/github_oauth_cb", githubLogin.getCallbackHandler())

	//https://github.com/dghubble/gologin/blob/master/google/login.go
	r.GET("/login_google", googleLogin.getLoginHandler())
	r.GET("/google_oauth_cb", googleLogin.getCallbackHandler())

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

	swaggerDir := filepath.Join(dir, "swagger")
	log.WithFields(log.Fields{"swaggerDir": swaggerDir}).Debug("http swaggerDir dir")
	r.StaticFS("/swagger", http.Dir(swaggerDir))

	apiV1 := r.Group("/api/v1")
	// https://goswagger.io/faq/faq_documenting.html#how-to-serve-swagger-ui-from-a-preexisting-web-app
	redoc := middleware.Redoc(middleware.RedocOpts{BasePath: "/api/v1", Path: "help", SpecURL: "/swagger/swagger.json", Title: "Hello"}, nil)
	apiV1.GET("/help", func(c *gin.Context) {
		redoc.ServeHTTP(c.Writer, c.Request)
	})

	apiV1.GET("/users", func(c *gin.Context) {
		id1 := int64(23)
		id2 := int64(33)
		c.JSON(http.StatusOK, []models.User{{ID: &id1, Name: "hongkai"}, {ID: &id2, Name: "mike"}})
	})
	apiV1.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		log.WithFields(log.Fields{"id": id}).Debug("getting user with id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, models.Error{Code: int64(http.StatusInternalServerError), Message: &msg})
			return
		}
		u := models.User{ID: &idInt64, Name: "hongkai"}
		c.JSON(http.StatusOK, u)
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}
