package cincinnati

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/yaml"

	"github.com/hongkailiu/test-go/pkg/testhelper"
)

func TestCincinnatiGraphData_LoadGraphData(t *testing.T) {
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

func TestCincinnatiGraphData_Shape(t *testing.T) {
	tests := []struct {
		name      string
		graphFile string
		gd        *CincinnatiGraphData
		edges     []Edge
		expectErr error
	}{
		{
			name:  "basic case.yaml",
			edges: []Edge{{1, 2}, {1, 0}, {2, 3}},
			gd: &CincinnatiGraphData{
				Channels: []Channel{
					{
						Name: "todo",
						Versions: []string{
							"4.10.0-rc.6",
							"4.10.0-rc.7",
							"4.10.0-rc.8",
							"4.10.2",
							"4.10.1",
							"4.10.1",
							"4.10.3",
							"4.10.11",
							"4.10.16",
							"4.10.18",
						},
					},
				},
				RemovedEdges: []RemovedEdge{{From: ".*", To: "4.10.0-rc.7"}},
				BlockedEdges: []BlockedEdge{
					{
						RemovedEdge:   RemovedEdge{From: "4.10.0-rc[.].*", To: "4.10.1"},
						Name:          "RiskA",
						MatchingRules: []MatchingRule{{Type: "Always"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", t.Name()+"_graph.json"))
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			var g Graph
			if err := json.Unmarshal(raw, &g); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			g.Edges = tt.edges
			actual, actualErr := tt.gd.Shape(context.TODO(), g)
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
