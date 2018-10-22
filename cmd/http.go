package main

import (
	"os"

	"github.com/hongkailiu/test-go/http"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true,})

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	http.Run()

}
