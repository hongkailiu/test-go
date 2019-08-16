package httpreverse

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hongkailiu/test-go/pkg/lib/util"
	log "github.com/sirupsen/logrus"
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

// Start http reverse server
func Start() {
	log.Info("Start http reverse server")

	c := loadConfig()
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: c.TargetScheme,
		Host:   c.TargetHost,
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.Port), Handler: proxy}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Info("http server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("error at server shutdown:", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("Server exited.")
	}

	log.Info("Server exited.")

}
