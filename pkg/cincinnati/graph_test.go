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

	g = g.EnsureEdges(edges)
	if len(g.Edges) == 0 {
		t.Errorf("Edges not match")
	}
}

func Test_RemoveEdges(t *testing.T) {
	tests := []struct {
		name         string
		g            Graph
		removedNodes []int
		expect       []Edge
	}{
		{
			name:         "remove 1",
			g:            Graph{Edges: []Edge{{1, 3}, {2, 3}, {0, 3}}},
			removedNodes: []int{1},
			expect:       []Edge{{1, 2}, {0, 2}},
		},
		{
			name:         "remove 1 2",
			g:            Graph{Edges: []Edge{{1, 3}, {2, 3}, {0, 3}}},
			removedNodes: []int{1, 2},
			expect:       []Edge{{0, 1}},
		},
		{
			name:         "remove 1 2 3",
			g:            Graph{Edges: []Edge{{1, 3}, {2, 3}, {0, 3}}},
			removedNodes: []int{1, 2, 3},
		},
		{
			name:         "another remove 1",
			g:            Graph{Edges: []Edge{{3, 1}, {2, 3}, {3, 2}}},
			removedNodes: []int{1},
			expect:       []Edge{{1, 2}, {2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.g.RemoveEdges(tt.removedNodes)
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("edges not match (-want +got):\n%s", diff)
			}
		})
	}
}
