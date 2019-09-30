package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

var (
	apiUrl = "http://api.openweathermap.org/data/2.5/weather"
	//http://api.openweathermap.org/data/2.5/weather?q=surrey,ca&appid=secret
)

type Service interface {
	GetWeather(city, country string) (Response, error)
	HandleResponse(r Response, names []string) error
}

type OpenWeatherMap struct {
	client *resty.Client
	appID  string
}

func NewOpenWeatherMap(appID string) *OpenWeatherMap {
	return &OpenWeatherMap{client: resty.New(), appID: appID}
}

func (w *OpenWeatherMap) GetWeather(city, country string) (Response, error) {
	resp, err := w.client.R().
		SetQueryParams(map[string]string{
			"q":     fmt.Sprintf("%s,%s", city, country),
			"appid": w.appID,
			"units": "metric",
		}).
		Get(apiUrl)
	if err != nil {
		return Response{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return Response{}, fmt.Errorf("status code is '%d', instead of 200", resp.StatusCode())
	}

	bytes := resp.Body()
	var response Response
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func getWriters(writerNames []string) ([]Writer, error) {
	var writers []Writer
	for _, writeName := range writerNames {
		switch writeName {
		case "logger":
			writers = append(writers, &Logger{logger: logrus.WithField("writer", "logger")})
		case "csv":
			writers = append(writers, &CSV{})
		case "yaml":
			writers = append(writers, &YAML{})
		default:
			return nil, fmt.Errorf("unknown write name: '%s'", writeName)
		}
	}
	return writers, nil
}

func (w *OpenWeatherMap) HandleResponse(r Response, writerNames []string) error {
	writers, err := getWriters(writerNames)
	if err != nil {
		return err
	}
	for _, w := range writers {
		if err := w.Write(r); err != nil {
			return err
		}
	}
	return nil
}
