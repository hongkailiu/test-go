package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProwJobAgent(t *testing.T) {
	testcases := []struct {
		description  string
		prowJobAgent ProwJobAgent
		expected     string
	}{
		{
			description:  "KubernetesAgent",
			prowJobAgent: KubernetesAgent,
			expected:     "kubernetes",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := string(tc.prowJobAgent)
			assert.Equal(t, tc.expected, result)
		})
	}

}
