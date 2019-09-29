package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
	apiUrl = "https://samples.openweathermap.org/data/2.5/weather"
	appID := "b6907d289e10d714a6e88b30761fae22"
	testcases := []struct {
		description      string
		city             string
		country          string
		expectedResponse Response
		expectedError    error
	}{
		{
			description: "London,uk",
			city:        "London",
			country:     "uk",
			expectedResponse: Response{
				ID:   2643743,
				Name: "London",
				CoOrd: CoOrd{
					Lat: float32(-0.13),
					Lon: float32(51.51),
				},
				Weather: []Weather{
					{
						ID:          300,
						Main:        "Drizzle",
						Description: "light intensity drizzle",
					},
				},
				Main: Main{
					Temp:     float32(280.32),
					Pressure: 1012,
					Humidity: 81,
					tempMin:  float32(0),
					tempMax:  float32(0),
				},
				Visibility: 10000,
				Wind: Wind{
					Speed: float32(4.1),
					Deg:   80,
				},
				Clouds:    map[string]int{"all": 90},
				Sys:       Sys{Country: "GB"},
				Condition: 200,
				Date:      NewJSONTime(time.Unix(1485789600, 0)),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			Service := NewOpenWeatherMap(appID)
			r, err := Service.getWeather(tc.city, tc.country)
			assert.Equal(t, tc.expectedResponse, r)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
