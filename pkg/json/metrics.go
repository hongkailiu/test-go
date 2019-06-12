package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	markerName string = "cluster_loader_marker"
)

type Metrics interface {
	printLog() error
}

type BaseMetrics struct {
	// To let the 3rd party know that this log entry is important
	// TODO set this up by config file
	Marker string `json:"marker"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type TestDuration struct {
	BaseMetrics
	StartTime    time.Time      `json:"startTime"`
	TestDuration time.Duration  `json:"testDuration"`
	Steps        []StepDuration `json:"Steps"`
}

type StepDuration interface {
}

type BaseStepDuration struct {
	Type      string        `json:"type"`
	StartTime time.Time     `json:"startTime"`
	TotalTime time.Duration `json:"totalTime"`
}

type TemplateStepDuration struct {
	BaseStepDuration
	RateDelay      time.Duration `json:"rateDelay"`
	RateDelayCount int           `json:"rateDelayCount"`
	StepPause      time.Duration `json:"stepPause"`
	StepPauseCount int           `json:"stepPauseCount"`
	SyncTime       time.Duration `json:"syncTime"`
}

type PodStepDuration struct {
	BaseStepDuration
	WaitPodsDurations []time.Duration `json:"waitPodsDurations"`
	RateDelay         time.Duration   `json:"rateDelay"`
	RateDelayCount    int             `json:"rateDelayCount"`
	StepPause         time.Duration   `json:"stepPause"`
	StepPauseCount    int             `json:"stepPauseCount"`
	SyncTime          time.Duration   `json:"syncTime"`
}

func (td TestDuration) printLog() error {
	b, err := json.Marshal(td)
	fmt.Println(string(b))
	return err
}

func LogMetrics(metrics []Metrics) error {
	for _, m := range metrics {
		err := m.printLog()
		if err != nil {
			return err
		}
	}
	return nil
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type MyTestDuration struct {
	TestDuration Duration         `json:"testDuration"`
	Steps        []MyStepDuration `json:"Steps"`
}

type MyStepDuration struct {
	StepDuration Duration         `json:"testDuration"`
	Steps        []MyStepDuration `json:"Steps"`
}
