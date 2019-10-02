package weather

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Writer interface {
	Write(r Record) error
}

type Logger struct {
	logger *logrus.Entry
}

func (w *Logger) Write(r Record) error {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	w.logger.Infof("%s", string(bytes))
	return nil
}

func (w *Logger) GetName() string {
	return "logger"
}

type CSV struct {
	OutputDir string
}

func (w *CSV) Write(r Record) (returnE error) {
	now := time.Now()
	rec, err := getRecord(now, r)
	if err != nil {
		return err
	}
	records := [][]string{
		rec,
	}

	f, err := os.OpenFile(path.Join(w.OutputDir, "weather.csv"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			returnE = err
		}
	}()

	csvWriter := csv.NewWriter(f)
	return csvWriter.WriteAll(records)
}

func getRecord(now time.Time, r Record) ([]string, error) {
	//{"first_name", "last_name", "username"}
	utcTime := r.Date.UTC()
	location, err := time.LoadLocation(r.TimeZone)
	if err != nil {
		return nil, err
	}
	timeInLocation := r.Date.In(location)
	logrus.WithField("utcTime.Format(http.TimeFormat)", utcTime.Format(http.TimeFormat)).
		WithField("r.TimeZone", r.TimeZone).
		WithField("location", fmt.Sprintf("%+v", location)).
		WithField("timeInLocation", timeInLocation).
		Debug("get record")
	return []string{
		now.UTC().Format(http.TimeFormat),
		r.Name,
		r.Sys.Country,
		timeInLocation.Format(http.TimeFormat),
		strconv.Itoa(timeInLocation.Year()),
		timeInLocation.Month().String(),
		strconv.Itoa(timeInLocation.Day()),
		strconv.Itoa(timeInLocation.Hour()),
		r.Weather[0].Main,
		r.Weather[0].Description,
		r.Weather[0].Icon,
		strconv.FormatFloat(r.Main.Temp, 'f', 1, 64),
	}, nil
}

type YAML struct {
}

func (w *YAML) Write(r Record) error {
	logrus.Warnf("not implemented yet! '%s'", "yaml")
	return nil
}
