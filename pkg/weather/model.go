package weather

import (
	"fmt"
	"strconv"
	"time"
)

//https://gist.github.com/uudashr/6b285cf0c44b0a7375d1b786967e1712
type JSONTime struct {
	time.Time
}

func NewJSONTime(t time.Time) JSONTime {
	return JSONTime{Time: t}
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", t.Time.Unix())
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil

}

type Response struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	CoOrd      CoOrd          `json:"coord"`
	Weather    []Weather      `json:"weather"`
	Main       Main           `json:"main"`
	Visibility int            `json:"visibility"`
	Wind       Wind           `json:"wind"`
	Clouds     map[string]int `json:"clouds"`
	Sys        Sys            `json:"sys"`
	Condition  int            `json:"cod"`
	Date       JSONTime       `json:"dt"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CoOrd struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Main struct {
	Temp     float64
	Pressure int
	Humidity int
	tempMin  float64
	tempMax  float64
}

type Wind struct {
	Speed float64
	Deg   float64
}

type Sys struct {
	Country string
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Record struct {
	Response
	TimeZone string
}

type TimezoneResponse struct {
	Status           string `json:"status"`
	Message          string `json:"message"`
	CountryCode      string `json:"countryCode"`
	CountryName      string `json:"countryName"`
	ZoneName         string `json:"zoneName"`
	Abbreviation     string `json:"abbreviation"`
	GmtOffset        int    `json:"gmtOffset"`
	Dst              string `json:"dst"`
	ZoneStart        int    `json:"zoneStart"`
	ZoneEnd          int    `json:"zoneEnd"`
	NextAbbreviation string `json:"nextAbbreviation"`
	Timestamp        int    `json:"timestamp"`
	Formatted        string `json:"formatted"`
}
