package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

	if err := c.Validate(); err != nil {
		logrus.WithError(err).WithField("o.configPath", o.configPath).Fatal("Failed to validate the config.")
	}

	logrus.Info("starting weather processing ...")

	service := weather.NewOpenWeatherMap(c.AppID, c.OutputDir)
	tzService := weather.NewTimeZoneDB(c.Key)

	var cities []weather.City
	var save = false
	for _, city := range c.Cities {
		r, err := service.GetWeather(city.Name, city.Country, false)
		if err != nil {
			logrus.WithError(err).WithField("city.Name", city.Name).WithField("city.Country", city.Country).Fatal("Failed to get weather.")
		}
		if len(city.TimeZone) == 0 {
			tzRes, err := tzService.GetTimeZone(r.CoOrd.Lat, r.CoOrd.Lon)
			if err != nil {
				logrus.WithError(err).WithField("r.CoOrd.Lon", r.CoOrd.Lon).
					WithField("r.CoOrd.Lat", r.CoOrd.Lat).
					WithField("response", fmt.Sprintf("%v", r)).
					Fatal("Failed to get timezone.")
			}
			save = true
			cities = append(cities, weather.City{Name: city.Name, Country: city.Country, TimeZone: tzRes.ZoneName})
			logrus.WithField("city.TimeZone", city.TimeZone).WithField("city.Name", city.Name).Debug("found the time zone of the city")
			logrus.Info("Sleeping 1 second due to TimeZoneDB limit")
			time.Sleep(1 * time.Second)
		} else {
			cities = append(cities, city)
		}
		record := weather.Record{Response: r, TimeZone: city.TimeZone}
		if err := service.HandleRecord(record, c.Writers); err != nil {
			logrus.WithError(err).WithField("record", fmt.Sprintf("%v", record)).Fatal("Failed to handle record.")
		}
	}
	if save {
		c.Cities = cities
		if err := c.Save(o.configPath); err != nil {
			logrus.WithError(err).Fatal("Failed to save config.")
		}
	}
}
