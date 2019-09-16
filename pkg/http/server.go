package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/securecookie"
	"github.com/hongkailiu/test-go/pkg/http/db"
	"github.com/hongkailiu/test-go/pkg/http/info"
	"github.com/hongkailiu/test-go/pkg/http/model"
	"github.com/hongkailiu/test-go/pkg/http/webhook"
	"github.com/hongkailiu/test-go/pkg/http/webhook/github"
	"github.com/hongkailiu/test-go/pkg/lib/util"
	"github.com/hongkailiu/test-go/pkg/random"
	"github.com/hongkailiu/test-go/pkg/swagger/swagger/models"
	cmdconfig "github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	appConfig *config
	log       *logrus.Logger
)

// PrometheusMiddleware intercepts all http requests and logging the path
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		log.WithFields(logrus.Fields{"c.Request.URL.Path": c.Request.URL.Path, "c.Request.Host": c.Request.Host, "c.Request.Header.Get(\"X-Forwarded-Host\")": c.Request.Header.Get("X-Forwarded-Host")}).
			Debug("prometheus middleware detected path visited")
		httpRequestsTotal.With(prometheus.Labels{"path": c.Request.URL.Path, "hostname": appConfig.hostname, "method": c.Request.Method}).Inc()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		log.WithFields(logrus.Fields{"c.Request.URL.Path": c.Request.URL.Path, "c.Writer.Status()": c.Writer.Status(), "latency": latency}).
			Debug("prometheus middleware calculated latency")
		httpRequestDuration.With(
			prometheus.Labels{"path": c.Request.URL.Path, "hostname": appConfig.hostname, "method": c.Request.Method, "status": strconv.Itoa(c.Writer.Status())}).
			Observe(latency.Seconds())
	}
}

// SetLog sets up log
func SetLog(l *logrus.Logger) {
	log = l
}

// Run starts the http server
func Run(hc *cmdconfig.HttpConfig) {
	log.WithFields(logrus.Fields{"hc.PProf": hc.PProf}).Info("Run with args")

	appConfig = loadConfig()

	oauthConfGitHub, err := setupOAuthConfig("github")
	if err != nil {
		log.Fatal("error at setupOAuthConfig:", err)
	}
	oauthConfGoogle, err := setupOAuthConfig("google")
	if err != nil {
		log.Fatal("error at setupOAuthConfig:", err)
	}

	githubLogin := login{oauthConfGitHub, gitHubUserProvider{}}
	googleLogin := login{oauthConfGoogle, googleUserProvider{}}

	appDBConfig := loadDBConfig()
	appDBString := appDBConfig.getDBString()

	//dbService *db.Service

	log.WithFields(logrus.Fields{"ClientID": oauthConfGitHub.ClientID, "ClientSecret": oauthConfGitHub.ClientSecret}).Debug("oauthConfGitHub")
	log.WithFields(logrus.Fields{"ClientID": oauthConfGoogle.ClientID, "ClientSecret": oauthConfGoogle.ClientSecret, "RedirectURL": oauthConfGoogle.RedirectURL}).Debug("oauthConfGoogle")

	log.WithFields(logrus.Fields{"appDBString": appDBString}).Debug("Using this db")
	appDB, err := db.OpenPostgres(appDBString)
	defer func() {
		err := appDB.Close()
		if err != nil {
			log.Fatal("error at appDB.Close():", err)
		}
	}()

	appDB.LogMode(true)

	if err != nil {
		log.Errorf("error occurred when db.OpenPostgres(appDBString): %s", err.Error())
		log.Warnf("running application without db")
	} else {
		db.Migrate(appDB)
	}
	dbService := db.New(appDB)

	prometheusRegister()
	generateRandomMetricsData()

	r := setupRouter(hc, githubLogin, googleLogin, dbService)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	//r.Run()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.WithFields(logrus.Fields{"port": port}).Debug("use port")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	util.ShutdownHTTPServer(srv, log)

	// for better test coverage
	if os.Getenv("unit_testing") == "true" {
		log.Info("Server not started because of unit testing.")
		return
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("Server exited.")
	}

	log.Info("Server exited.")
}

func setupOAuthConfig(cloud string) (*oauth2.Config, error) {
	switch cloud {
	case "github":
		// You must register the app at https://github.com/settings/applications
		// Set callback to http://127.0.0.1:7000/github_oauth_cb
		// Set ClientId and ClientSecret to
		return &oauth2.Config{
			ClientID:     appConfig.ghClientID,
			ClientSecret: appConfig.ghClientSecret,
			// select level of access you want from https://developer.github.com/v3/oauth/#scopes
			Scopes:   []string{"read:user", "user:email"},
			Endpoint: githuboauth.Endpoint,
		}, nil
	case "google":
		//https://console.developers.google.com/apis/dashboard
		return &oauth2.Config{
			ClientID:     appConfig.ggClientID,
			ClientSecret: appConfig.ggClientSecret,
			RedirectURL:  appConfig.ggRedirectURL,
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		}, nil
	default:
		return nil, fmt.Errorf("do not support oauth for provider %s", cloud)
	}

}

func generateRandomMetricsData() {
	go func() {
		for {
			n := random.GetRandom(1000)
			log.WithFields(logrus.Fields{"n": n}).Debug("generated random number")
			randomNumber.With(prometheus.Labels{"key": "value", "hostname": appConfig.hostname}).Set(float64(n))
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			n := random.GetRandom(50)
			log.WithFields(logrus.Fields{
				"n": n,
			}).Debug("generated fake storageOperationMetric")
			storageOperationMetric.With(prometheus.Labels{"volume_plugin": "hliu.ca/aws-ebs",
				"operation_name": "volume_provision", "hostname": appConfig.hostname}).Observe(float64(n))
			time.Sleep(100 * time.Second)
		}
	}()
}

func setupRouter(hc *cmdconfig.HttpConfig, githubLogin, googleLogin login, dbService db.ServiceI) *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(PrometheusMiddleware())

	// sessions and cookies are from github.com/gin-contrib
	// which uses implementation from github.com/gorilla/
	log.WithFields(logrus.Fields{"sessionKey": appConfig.sessionKey}).Info("using session key")
	store := cookie.NewStore(appConfig.sessionKey)
	r.Use(sessions.Sessions("my_session", store))

	r.GET("/", func(c *gin.Context) {
		infoP := info.GetInfo(hc.Version)
		contentType := c.ContentType()
		log.WithFields(logrus.Fields{"c.ContentType()": contentType}).Debug("root path")
		if strings.Contains(contentType, "yaml") {
			c.YAML(http.StatusOK, *infoP)
		} else {
			c.JSON(http.StatusOK, *infoP)
		}
	})

	// https://developer.github.com/webhooks/
	r.POST("/webhook", func(c *gin.Context) {
		// TODO ShouldBindHeader will be available in the next gin release
		//if err := c.Request.ShouldBindHeader(&h); err != nil {
		//	c.JSON(200, err)
		//}
		h := github.WebhookHeaders{
			EventType:     c.GetHeader("X-GitHub-Event"),
			GUID:          c.GetHeader("X-GitHub-Delivery"),
			HMACHexDigest: c.GetHeader("X-Hub-Signature"),
			ContentType:   c.GetHeader("content-type"),
		}
		log.WithFields(logrus.Fields{"h": fmt.Sprintf("%+v", h)}).Debug("webhook headers")
		payload, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, webhook.Response{Message: "cannot read request body"})
			return
		}
		status, msg := github.Handle(payload, h, log)
		c.JSON(status, webhook.Response{Message: msg})
	})

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	staticDir := filepath.Join(dir, "static")
	log.WithFields(logrus.Fields{"staticDir": staticDir}).Debug("http staticDir dir")
	r.StaticFS("/console", http.Dir(staticDir))

	//https://blog.kowalczyk.info/article/f/accessing-github-api-from-go.html
	r.GET("/login_github", githubLogin.getLoginHandler())
	r.GET("/github_oauth_cb", githubLogin.getCallbackHandler())

	//https://github.com/dghubble/gologin/blob/master/google/login.go
	r.GET("/login_google", googleLogin.getLoginHandler())
	r.GET("/google_oauth_cb", googleLogin.getCallbackHandler())

	r.GET("/whoami", func(c *gin.Context) {
		u := getKeyInSession(c, "username")
		username := ""
		if u != nil {
			username = *u
		}
		c.JSON(200, gin.H{"username": username})
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		session.Delete("username")
		session.Delete("localID")
		err := session.Save()
		if err != nil {
			log.WithError(err).Errorf("error when user %s logout.", username)
		}
		c.Redirect(http.StatusTemporaryRedirect, "/console")
	})

	swaggerDir := filepath.Join(dir, "swagger")
	log.WithFields(logrus.Fields{"swaggerDir": swaggerDir}).Debug("http swaggerDir dir")
	r.StaticFS("/swagger", http.Dir(swaggerDir))

	apiV1 := r.Group("/api/v1")
	apiV1.Use(AuthorizationMiddleware())
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
		log.WithFields(logrus.Fields{"id": id}).Debug("getting user with id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, models.Error{Code: int64(http.StatusInternalServerError), Message: &msg})
			return
		}
		u := models.User{ID: &idInt64, Name: "hongkai"}
		c.JSON(http.StatusOK, u)
	})

	apiV1.GET("/cities", func(c *gin.Context) {
		cities, err := dbService.GetCities(10, 0)
		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, models.Error{Code: int64(http.StatusInternalServerError), Message: &msg})
			return
		}
		if cities == nil {
			cities = &[]model.City{}
		}
		c.JSON(http.StatusOK, *cities)
	})

	r.GET("/token", AuthenticationMiddleware(), func(c *gin.Context) {
		localID := getKeyInSession(c, "localID")
		tokenString, err := getToken(*localID, appConfig.sessionKey)

		if err != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, models.Error{Code: int64(http.StatusInternalServerError), Message: &msg})
			return
		}
		log.WithFields(logrus.Fields{"tokenString": tokenString}).Debug("generated token")
		c.JSON(200, gin.H{"token": tokenString})
	})

	// https://github.com/gin-contrib/pprof
	// default is "debug/pprof"
	if hc.PProf {
		pprof.Register(r)
	}

	return r
}

func getKeyInSession(c *gin.Context, key string) *string {
	if value, exists := c.Get(sessions.DefaultKey); exists {
		session := value.(sessions.Session)
		v := session.Get(key)
		if v == nil {
			return nil
		}
		username := v.(string)
		return &username
	}
	return nil
}

// GetSecret returns a secret
func GetSecret(length int) []byte {
	return securecookie.GenerateRandomKey(length)
}
