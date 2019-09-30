package weather

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestLoadConfig(t *testing.T) {
	testcases := []struct {
		description    string
		path           string
		expectedConfig Config
		expectedError  error
	}{
		{
			description: "test_file",
			path:        "test_files/test_config.yaml",
			expectedConfig: Config{
				AppID: "abc",
				Cities: []City{
					{
						Name:    "vancouver",
						Country: "ca",
					},
					{
						Name:    "toronto",
						Country: "ca",
					},
					{
						Name:    "montreal",
						Country: "ca",
					},
				},
				Writers:   []string{"logger", "csv", "yaml"},
				OutputDir: "/tmp",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			c, err := LoadConfig(tc.path)
			assert.Equal(t, tc.expectedConfig, c)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
