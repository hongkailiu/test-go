package cincinnati

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/go-cmp/cmp"
)

func TestNode_getPrevious(t *testing.T) {
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

func TestGraph_RemoveEdges(t *testing.T) {
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

func TestGraph_getTag(t *testing.T) {
	tests := []struct {
		name        string
		g           Graph
		version     string
		suffix      ArchTagSuffix
		expectTag   string
		expectIndex int
	}{
		{
			name: "4.2.11",
			g: Graph{Nodes: []Node{
				{
					Version: semver.MustParse("4.2.11"),
					Tag:     "4.2.11",
				},
			}},
			version:   "4.2.11",
			suffix:    ArchTagSuffixAMD64,
			expectTag: "4.2.11",
		},
		{
			name: "4.3.11-x86_64",
			g: Graph{Nodes: []Node{
				{
					Version: semver.MustParse("4.3.11"),
					Tag:     "4.3.11-x86_64",
				},
			}},
			version:   "4.3.11",
			suffix:    ArchTagSuffixAMD64,
			expectTag: "4.3.11-x86_64",
		},
		{
			name: "4.23.11-multi",
			g: Graph{Nodes: []Node{
				{
					Version: semver.MustParse("4.23.11"),
					Tag:     "4.23.11-multi",
				},
			}},
			version:   "4.23.11",
			suffix:    ArchTagSuffixMULTI,
			expectTag: "4.23.11-multi",
		},
		{
			name: "4.23.11-multi",
			g: Graph{Nodes: []Node{
				{
					Version: semver.MustParse("4.23.11"),
					Tag:     "4.23.11-multi",
				},
			}},
			version:     "4.23.11",
			suffix:      ArchTagSuffixAMD64,
			expectIndex: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTag, actualIndex := tt.g.getTagAndIndex(tt.version, tt.suffix)

			if diff := cmp.Diff(tt.expectIndex, actualIndex); diff != "" {
				t.Errorf("index not match (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.expectTag, actualTag); diff != "" {
				t.Errorf("tag not match (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGraph_compatible(t *testing.T) {
	tests := []struct {
		name string
		g    Graph
	}{
		{
			name: "empty graph.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.g.compatible()

			data, err := json.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}

			cupaloy.New(
				cupaloy.SnapshotSubdirectory("testdata/.snapshots"),
			).SnapshotT(t, data)
		})
	}
}

func TestNode_getFrom(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		version string
		graph   Graph
		expect  int
	}{
		{
			name:   "empty graph and node",
			expect: -1,
		},
		{
			name:    "no arch suffix",
			version: "1.0.0",
			graph: Graph{
				Nodes: []Node{
					{
						Version: semver.MustParse("0.0.1"),
						Tag:     "0.0.1",
					},
					{
						Version: semver.MustParse("1.0.0"),
						Tag:     "1.0.0",
					},
				},
			},
			node: Node{
				Version: semver.MustParse("1.0.0"),
				Tag:     "1.0.0",
			},
			expect: 1,
		},
		{
			name:    "4.2.11-s390x",
			version: "4.2.11-s390x",
			graph: Graph{
				Nodes: []Node{
					{
						Version: semver.MustParse("0.0.1"),
						Tag:     "0.0.1",
					},
					{
						Version: semver.MustParse("4.2.11-s390x"),
						Tag:     "4.2.11-s390x",
					},
				},
			},
			node: Node{
				Version: semver.MustParse("4.2.11-s390x"),
				Tag:     "4.2.11-s390x",
			},
			expect: 1,
		},
		{
			name:    "4.13.11-s390x",
			version: "4.13.11",
			graph: Graph{
				Nodes: []Node{
					{
						Version: semver.MustParse("0.0.1"),
						Tag:     "0.0.1",
					},
					{
						Version: semver.MustParse("4.13.11"),
						Tag:     "4.13.11-s390x",
					},
				},
			},
			node: Node{
				Version: semver.MustParse("4.13.11"),
				Tag:     "4.13.11-s390x",
			},
			expect: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.getFrom(tt.version, tt.graph)
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("index not match (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCincinnatiGraphData_BecomeConditional(t *testing.T) {
	tests := []struct {
		name   string
		gd     CincinnatiGraphData
		from   Node
		to     string
		expect ConditionalEdge
	}{
		{
			name: "arch suffix works",
			gd: CincinnatiGraphData{
				BlockedEdges: []BlockedEdge{
					{
						RemovedEdge: RemovedEdge{
							From: `^4[.](17[.](2[0-8]|[1]?[0-9])|18[.](1[01]|[0-9]))[+].*$`,
							To:   "4.18.12",
						},
						MatchingRules: []MatchingRule{
							{
								Type: "Always",
							},
						},
					},
				},
			},
			from: Node{
				Tag:     "4.18.10-x86_64",
				Version: semver.MustParse("4.18.10"),
			},
			to: "4.18.12",
			expect: ConditionalEdge{
				RisksKey: "[0]",
				Edges: []ConditionalUpdate{
					{
						From: "4.18.10",
						To:   "4.18.12",
					},
				},
				Risks: []ConditionalUpdateRisk{
					{
						MatchingRules: []MatchingRule{
							{
								Type: "Always",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.gd.BecomeConditional(tt.from, tt.to)
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("conditional edge not match (-want +got):\n%s", diff)
			}
		})
	}
}
