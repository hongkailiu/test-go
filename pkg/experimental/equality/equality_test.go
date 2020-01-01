package equality

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestEqualSet(t *testing.T) {
	testCases := []struct {
		name string
		arg1 sets.String
		arg2 sets.String
	}{
		{
			name: "elem ordering in set does not matter",
			arg1: sets.NewString("1", "2"),
			arg2: sets.NewString("2", "1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.arg1, tc.arg2) {
				t.Errorf("unexpected non equal (reflect): %s", cmp.Diff(tc.arg1, tc.arg2))
			}
			if !equality.Semantic.DeepEqual(tc.arg1, tc.arg2) {
				t.Errorf("unexpected non equal (equality): %s", cmp.Diff(tc.arg1, tc.arg2))
			}
		})
	}
}
