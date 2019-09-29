package weather

import (
	"fmt"
	"strconv"
	"time"
)

type JSONTime time.Time

func NewJSONTime(t time.Time) JSONTime {
	z := JSONTime(t)
	return z
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", time.Time(*t).Unix())
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = NewJSONTime(time.Unix(i, 0))
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
}

type CoOrd struct {
	Lat float32 `json:"lon"`
	Lon float32 `json:"lat"`
}

type Main struct {
	Temp     float32
	Pressure int
	Humidity int
	tempMin  float32
	tempMax  float32
}

type Wind struct {
	Speed float32
	Deg   int
}

type Sys struct {
	Country string
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
