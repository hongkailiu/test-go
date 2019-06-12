package json

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypeDuration(t *testing.T) {
	objectUnderTest := 3*time.Second + 5*time.Minute
	actual := fmt.Sprintf("%s", objectUnderTest)
	assert.Equal(t, "5m3s", actual, "they should be equal")
}

func TestTypeDuration1(t *testing.T) {
	td := TestDuration{
		TestDuration: 23*time.Second + 6*time.Minute,
		Steps: []StepDuration{
			TemplateStepDuration{
				RateDelay: 2 * time.Second,
				SyncTime:  3 * time.Second,
			},
			PodStepDuration{},
		},
	}

	bolB, err := json.Marshal(td)
	assert.Nil(t, err)

	actual := string(bolB)

	fmt.Println("=" + actual + "=")

	expect := `{"marker":"","name":"","type":"","startTime":"0001-01-01T00:00:00Z","testDuration":383000000000,"Steps":[{"type":"","startTime":"0001-01-01T00:00:00Z","totalTime":0,"rateDelay":2000000000,"rateDelayCount":0,"stepPause":0,"stepPauseCount":0,"syncTime":3000000000},{"type":"","startTime":"0001-01-01T00:00:00Z","totalTime":0,"waitPodsDurations":null,"rateDelay":0,"rateDelayCount":0,"stepPause":0,"stepPauseCount":0,"syncTime":0}]}`

	assert.Equal(t, expect, actual, "they should be equal")
}
