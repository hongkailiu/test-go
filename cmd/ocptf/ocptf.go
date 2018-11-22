package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hongkailiu/test-go/pkg/ocptf"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("ocp-tf", "A script to read terraform tfstate file and output inventory.")

	_    = app.Version(ocptf.VERSION).HelpFlag.Short('h')
	list = app.Flag("list", "List dynamic inventory.").Bool()
	host = app.Flag("host", "List dynamic inventory for host.").String()
)

func init() {
	if strings.ToLower(os.Getenv("ocp_tf_debug")) == "true" {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	log.WithFields(log.Fields{"list": *list, "host": *host}).Debug("args")

	if !*list && len(*host) == 0 {
		fmt.Fprintf(os.Stderr, "illegal args try '%s -h'\n", os.Args[0])
		os.Exit(1)
		//app.FatalUsage("illegal args try '%s -h'\n", os.Args[0])
	}

	path, err := getTerraformTFStateFile()
	log.WithFields(log.Fields{"path": *path}).Debug("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	dynamic := true
	if strings.ToLower(os.Getenv("static_inventory")) == "true" {
		dynamic = false
	}
	log.WithFields(log.Fields{"dynamic": dynamic}).Debug("")
	if *list {
		err = ocptf.DoList(*path, dynamic)
	} else {
		err = ocptf.DoHost(*path, *host, dynamic)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

}
func getTerraformTFStateFile() (*string, error) {
	path := os.Getenv("terraform_tf_state_file")
	if len(path) != 0 {
		return &path, nil
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	path = filepath.Join(dir, "terraform.tfstate")
	return &path, nil
}