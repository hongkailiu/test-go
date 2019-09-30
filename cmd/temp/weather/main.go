package main

import (
	"flag"
	"fmt"
	"os"

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

	c, err := weather.LoadConfig(o.configPath)
	if err != nil {
		logrus.WithError(err).WithField("o.configPath", o.configPath).Fatal("Failed to load the config.")
	}

	logrus.Info("starting weather processing ...")

	service := weather.NewOpenWeatherMap(c.AppID)

	for _, city := range c.Cities {
		r, err := service.GetWeather(city.Name, city.Country)
		if err != nil {
			logrus.WithError(err).WithField("city.Name", city.Name).WithField("city.Country", city.Country).Fatal("Failed to get weather.")
		}
		if err := service.HandleResponse(r, c.Writers); err != nil {
			logrus.WithError(err).WithField("response", fmt.Sprintf("%v", r)).Fatal("Failed to handle response.")
		}
	}

}
