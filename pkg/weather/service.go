package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

const (
	apiUrl       = "http://api.openweathermap.org/data/2.5/weather"
	apiSampleUrl = "https://samples.openweathermap.org/data/2.5/weather"
	//http://api.openweathermap.org/data/2.5/weather?q=surrey,ca&appid=secret

	apiTimeZoneUrl = "http://api.timezonedb.com/v2.1/get-time-zone"
	//http://api.timezonedb.com/v2.1/get-time-zone?key=value&format=json&by=position&lat=43.65&lng=-79.39
)

type Service interface {
	GetWeather(city, country string, sample bool) (Response, error)
	HandleRecord(r Record, names []string) error
}

type OpenWeatherMap struct {
	client       *resty.Client
	appID        string
	outputFolder string
}

func NewOpenWeatherMap(appID, outputFolder string) *OpenWeatherMap {
	return &OpenWeatherMap{client: resty.New(), appID: appID, outputFolder: outputFolder}
}

func (w *OpenWeatherMap) GetWeather(city, country string, sample bool) (Response, error) {
	url := apiUrl
	params := map[string]string{
		"q":     fmt.Sprintf("%s,%s", city, country),
		"appid": w.appID,
		"units": "metric",
	}
	if sample {
		url = apiSampleUrl
		delete(params, "units")
	}
	resp, err := w.client.R().
		SetQueryParams(params).
		Get(url)
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

func getWriters(writerNames []string, outputDir string) ([]Writer, error) {
	var writers []Writer
	for _, writeName := range writerNames {
		switch writeName {
		case "logger":
			writers = append(writers, &Logger{logger: logrus.WithField("writer", "logger")})
		case "csv":
			writers = append(writers, &CSV{OutputDir: outputDir})
		case "yaml":
			writers = append(writers, &YAML{})
		default:
			return nil, fmt.Errorf("unknown write name: '%s'", writeName)
		}
	}
	return writers, nil
}

func (w *OpenWeatherMap) HandleRecord(r Record, writerNames []string) error {
	writers, err := getWriters(writerNames, w.outputFolder)
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

type TimeZoneService interface {
	GetWeather(city, country string, sample bool) (TimezoneResponse, error)
}

func NewTimeZoneDB(key string) *TimeZoneDB {
	return &TimeZoneDB{client: resty.New(), key: key}
}

type TimeZoneDB struct {
	client *resty.Client
	key    string
}

func (tz *TimeZoneDB) GetTimeZone(lat, lng float64) (TimezoneResponse, error) {
	url := apiTimeZoneUrl
	params := map[string]string{
		"lat":    strconv.FormatFloat(lat, 'f', 2, 64),
		"lng":    strconv.FormatFloat(lng, 'f', 2, 64),
		"key":    tz.key,
		"format": "json",
		"by":     "position",
	}
	resp, err := tz.client.R().
		SetQueryParams(params).
		Get(url)
	if err != nil {
		return TimezoneResponse{}, err
	}

	bytes := resp.Body()
	logrus.WithField("string(bytes)", string(bytes)).Debugf("response from '%s'", apiTimeZoneUrl)

	if resp.StatusCode() != http.StatusOK {
		return TimezoneResponse{}, fmt.Errorf("status code from '%s' is '%d', instead of 200", apiTimeZoneUrl, resp.StatusCode())
	}

	var response TimezoneResponse
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return TimezoneResponse{}, err
	}
	return response, nil
}
