package weather

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Writer interface {
	Write(response Response) error
}

type Logger struct {
	logger *logrus.Entry
}

func (w *Logger) Write(response Response) error {
	bytes, err := json.Marshal(response)
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

func (w *CSV) Write(response Response) (returnE error) {
	now := time.Now()
	records := [][]string{
		getRecord(now, response),
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

func getRecord(now time.Time, r Response) []string {
	//{"first_name", "last_name", "username"}
	utcTime := r.Date.UTC()
	logrus.WithField("utcTime.Format(http.TimeFormat)", utcTime.Format(http.TimeFormat)).Debug("get record")
	return []string{
		now.UTC().Format(http.TimeFormat),
		r.Name,
		r.Sys.Country,
		utcTime.Format(http.TimeFormat),
		strconv.Itoa(utcTime.Year()),
		utcTime.Month().String(),
		strconv.Itoa(utcTime.Day()),
		strconv.Itoa(utcTime.Hour()),
		r.Weather[0].Main,
		r.Weather[0].Description,
		r.Weather[0].Icon,
		strconv.FormatFloat(r.Main.Temp, 'f', 1, 64),
	}
}

type YAML struct {
}

func (w *YAML) Write(response Response) error {
	logrus.Warnf("not implemented yet! '%s'", "yaml")
	return nil
}
