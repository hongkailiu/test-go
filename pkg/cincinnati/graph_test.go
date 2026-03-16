package cincinnati

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_getPrevious(t *testing.T) {
	var g Graph
	raw, err := os.ReadFile(filepath.Join("testdata", t.Name()+"_graph.json"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if err := json.Unmarshal(raw, &g); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	var found Node
	for _, node := range g.Nodes {
		if node.Version.String() == "4.10.1" {
			found = node
		}
	}
	if found.Tag == "" {
		t.Errorf("Failed to find tag")
	}
	edges := found.getPrevious(g)
	if diff := cmp.Diff([]Edge{{1, 3}, {2, 3}, {0, 3}}, edges); diff != "" {
		t.Errorf("edges not match (-want +got):\n%s", diff)
	}

	EnsureEdges(&g, edges)
	if len(g.Edges) == 0 {
		t.Errorf("Edges not match")
	}
}
