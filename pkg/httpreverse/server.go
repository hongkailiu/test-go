package httpreverse

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hongkailiu/test-go/pkg/lib/util"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

type config struct {
	Port string

	TargetScheme string
	TargetHost   string
}

func loadConfig() config {
	c := config{}
	c.Port = util.Getenv("port", "8888")
	c.TargetScheme = util.Getenv("target_scheme", "http")
	c.TargetHost = util.Getenv("target_host", "localhost:8080")
	return c
}

func setupReverseProxy(targetURL *url.URL) *httputil.ReverseProxy {

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	director := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		director(req)
		//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Host
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Host = req.URL.Host
	}
	return reverseProxy
}

// SetLog sets up log
func SetLog(l *logrus.Logger) {
	log = l
}

// Start http reverse server
func Start() {
	log.Info("Start http reverse server")

	c := loadConfig()
	log.WithField("config", fmt.Sprintf("%+v", c)).Info("config loaded")

	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.Port), Handler:
	setupReverseProxy(&url.URL{Scheme: c.TargetScheme, Host: c.TargetHost,})}

	util.ShutdownHTTPServer(srv, log)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("Server exited.")
	}

	log.Info("Server exited.")

}
