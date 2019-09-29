package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/resty.v1"
)

var (
	apiUrl = "http://api.openweathermap.org/data/2.5/weather"
	//http://api.openweathermap.org/data/2.5/weather?q=surrey,ca&appid=secret
)

type Service interface {
	getWeather(city, country string) (Response, error)
}

type OpenWeatherMap struct {
	client *resty.Client
	appID  string
}

func NewOpenWeatherMap(appID string) *OpenWeatherMap {
	return &OpenWeatherMap{client: resty.New(), appID: appID}
}

func (w *OpenWeatherMap) getWeather(city, country string) (Response, error) {
	resp, err := w.client.R().
		SetQueryParams(map[string]string{
			"q":     fmt.Sprintf("%s,%s", city, country),
			"appid": w.appID,
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
