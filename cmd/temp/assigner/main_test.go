package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/hongkailiu/test-go/pkg/assigner"
)

func TestDetermine(t *testing.T) {
	testCases := []struct {
		name          string
		now           time.Time
		config        assigner.Config
		expected      []string
		expectedError error
	}{
		{
			name:          "Nothing to schedule with empty config",
			now:           time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			expectedError: fmt.Errorf("not found"),
		},
		{
			name: "Nothing to schedule yet",
			now:  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			config: assigner.Config{
				GroupName: "",
				ScheduledActions: []assigner.Action{
					{
						At:      time.Date(2009, 12, 17, 20, 34, 58, 651387237, time.UTC),
						Members: nil,
					},
				},
			},
			expectedError: fmt.Errorf("not found"),
		},
		{
			name: "Nothing to schedule yet",
			now:  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			config: assigner.Config{
				GroupName: "",
				ScheduledActions: []assigner.Action{
					{
						At:      time.Date(2009, 10, 17, 20, 34, 58, 651387237, time.UTC),
						Members: []string{"a"},
					},
				},
			},
			expected: []string{"a"},
		},
		{
			name: "Nothing to schedule yet",
			now:  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			config: assigner.Config{
				GroupName: "",
				ScheduledActions: []assigner.Action{
					{
						At:      time.Date(2009, 12, 17, 20, 34, 58, 651387237, time.UTC),
						Members: []string{"c"},
					},
					{
						At:      time.Date(2009, 10, 17, 20, 34, 58, 651387237, time.UTC),
						Members: []string{"b"},
					},
					{
						At:      time.Date(2009, 9, 17, 20, 34, 58, 651387237, time.UTC),
						Members: []string{"a"},
					},
				},
			},
			expected: []string{"b"},
		},
	}

	for _, tc := range testCases {
		actual, err := determine(tc.now, tc.config)
		equal(t, tc.expected, actual)
		equalError(t, tc.expectedError, err)
	}
}

func equalError(t *testing.T, expected, actual error) {
	if expected != nil && actual == nil || expected == nil && actual != nil {
		t.Errorf("expecting error \"%v\", got \"%v\"", expected, actual)
	}
	if expected != nil && actual != nil && expected.Error() != actual.Error() {
		t.Errorf("expecting error msg %q, got %q", expected.Error(), actual.Error())
	}
}

func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("actual differs from expected:\n%s", cmp.Diff(expected, actual))
	}
}
