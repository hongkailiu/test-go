package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
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
					Lat: 51.51,
					Lon: -0.13,
				},
				Weather: []Weather{
					{
						ID:          300,
						Main:        "Drizzle",
						Description: "light intensity drizzle",
						Icon:        "09d",
					},
				},
				Main: Main{
					Temp:     280.32,
					Pressure: 1012,
					Humidity: 81,
					tempMin:  float64(0),
					tempMax:  float64(0),
				},
				Visibility: 10000,
				Wind: Wind{
					Speed: 4.1,
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
			service := NewOpenWeatherMap(appID, "")
			r, err := service.GetWeather(tc.city, tc.country, true)
			assert.Equal(t, tc.expectedResponse, r)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
