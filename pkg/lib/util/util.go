package util

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GetRandomInt return a random integer
func GetRandomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

// PanicIfError panics if error is not nil
func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Getenv returns value of env. var if it is defined; defaultValue otherwise.
func Getenv(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}

// HomeDir returns home dir
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

type log interface {
	Println(v ...interface{})
	Fatal(v ...interface{})
}

// ShutdownHTTPServer shuts down the http server
func ShutdownHTTPServer(server *http.Server, l log) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		l.Println("http server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			l.Fatal("error at server shutdown (ShutdownHTTPServer):", err)
		}
	}()
}
