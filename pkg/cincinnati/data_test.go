package cincinnati

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/yaml"

	"github.com/hongkailiu/test-go/pkg/testhelper"
)

func Test_LoadGraphData(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		expect    *CincinnatiGraphData
		expectErr error
	}{
		{
			name: "mock dir.yaml",
			dir:  "../../mock/graph-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := LoadGraphData(tt.dir)
			if diff := cmp.Diff(tt.expectErr, actualErr, testhelper.ErrorTransformer); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}

			data, err := yaml.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}

			cupaloy.New(
				cupaloy.SnapshotSubdirectory("testdata/.snapshots"),
			).SnapshotT(t, data)
		})
	}
}
