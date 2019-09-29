package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/hongkailiu/test-go/pkg/weather"
)

type options struct {
	configPath string
	debugMode  bool
}

func parseOptions() options {
	var o options
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&o.configPath, "config", "config.yaml", "The path to load config file.")
	fs.BoolVar(&o.debugMode, "debug-mode", false, "Enable the DEBUG level of logs if true.")
	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.WithError(err).Errorf("cannot parse args: '%s'", os.Args[1:])
	}
	return o
}

func main() {
	o := parseOptions()

	if o.debugMode {
		logrus.Info("debug mode is enabled")
		logrus.SetLevel(logrus.DebugLevel)
	}

	_, err := weather.LoadConfig(o.configPath)
	if err != nil {
		logrus.WithError(err).WithField("o.configPath", o.configPath).Fatal("Failed to load the config.")
	}

	logrus.Info("starting weather processing ...")

	logrus.Info("configure cron jobs ...")
	cron := cron.New()
	err = cron.AddFunc("*/10 * * * * *", func() {
		now := time.Now()
		logrus.WithField("now", now.Format(time.RFC3339)).Debug("Every 10 seconds ... ")
		//TODO
	})
	if err != nil {
		logrus.WithError(err).Fatal("error occurred when the add cron job")
	}
	cron.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		logrus.Println("weather shutting down...")
		cron.Stop()
	}()

	select {}
}
