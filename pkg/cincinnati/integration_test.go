package cincinnati

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/hongkailiu/test-go/pkg/util"
)

// TODO: compare result with production

func (g Graph) Equal(g1 Graph) bool {
	return g.IsSuperGraphOf(g1) && g1.IsSuperGraphOf(g)
}

func (g Graph) IsSuperGraphOf(g1 Graph) bool {
	return g.IsNodesSupersetOf(g1) && g.IsEdgesSupersetOf(g1) && g.IsConditionalEdgesSupersetOf(g1)
}

func (g Graph) IsNodesSupersetOf(g1 Graph) bool {
	for i, n := range g1.Nodes {
		if g.FindNode(n) == -1 {
			logrus.WithField("nodeVersion", n.Version).WithField("nodeImage", n.Image).WithField("tag", n.Tag).
				WithField("index", i).Error("node is not in graph")
			return false
		}
	}
	return true
}

func (g Graph) IsEdgesSupersetOf(g1 Graph) bool {
	for i, edge := range g1.Edges {
		var from, to int
		for j, n := range []Node{g1.Nodes[edge[0]], g1.Nodes[edge[1]]} {
			index := g.FindNode(g1.Nodes[edge[0]])
			if index == -1 {
				logrus.WithField("nodeVersion", n.Version).WithField("nodeImage", n.Image).WithField("tag", n.Tag).
					WithField("index", j).Error("node is not in graph")
				return false
			}
			if j == 0 {
				from = index
			} else {
				to = index
			}
		}
		if g.FindEdge(Edge{from, to}) == -1 {
			logrus.WithField("index", i).WithField("edge", edge).Error("edge is not in graph")
			return false
		}
	}
	return true
}

func (n Node) Equal(n1 Node) bool {
	return n.Image == n1.Image
}

func (g Graph) FindNode(node Node) int {
	for i, n := range g.Nodes {
		if n.Equal(node) {
			return i
		}
	}
	return -1
}

func (g Graph) IsConditionalEdgesSupersetOf(g1 Graph) bool {
	for _, ce1 := range g1.ConditionalEdges {
		for _, e1 := range ce1.Edges {
			var found bool
			for _, ce := range g.ConditionalEdges {
				for _, e := range ce.Edges {
					if e1.From == e.From && e1.To == e.To {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				logrus.WithField("conditionalUpdateEdge", e1).Error("conditional update edge is not in graph")
				return false
			}
		}
	}
	return true
}

func TestIntegration_dummy(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "1" {
		t.Skip("integration tests skipped unless TEST_INTEGRATION=1")
	}
	data, err := util.ReadFileMaybeGZIP(filepath.Join("../../data", "graph.json.gz"))
	if err != nil {
		t.Fatal(err)
	}
	graph := Graph{}
	err = json.Unmarshal(data, &graph)
	if err != nil {
		t.Fatal(err)
	}

	data, err = util.ReadFileMaybeGZIP(filepath.Join("../../data", "production_stable-4.10_amd64.json.gz"))
	if err != nil {
		t.Fatal(err)
	}
	production := Graph{}
	err = json.Unmarshal(data, &production)
	if err != nil {
		t.Fatal(err)
	}
	if !production.IsSuperGraphOf(graph) {
		t.Fatal("production is not a super-graph of graph")
	}
}
