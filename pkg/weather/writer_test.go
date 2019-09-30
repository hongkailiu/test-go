package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRecord(t *testing.T) {
	now := time.Unix(1569878592, 0)
	testcases := []struct {
		description string
		response    Response
		expected    []string
	}{
		{
			description: "London,uk",
			response: Response{
				ID:   2643743,
				Name: "London",
				CoOrd: CoOrd{
					Lat: -0.13,
					Lon: 51.51,
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
			expected: []string{"Mon, 30 Sep 2019 21:23:12 GMT", "Mon, 30 Jan 2017 15:20:00 GMT", "2017", "January", "30", "15", "Drizzle", "light intensity drizzle", "09d", "280.3"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := getRecord(now, tc.response)
			assert.Equal(t, tc.expected, result)
		})
	}

}
