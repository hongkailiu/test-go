package cincinnati

import (
	"os"
	"testing"
)

// TODO: compare result with production

func (g Graph) Equal(g1 Graph) bool {
	return g.IsSuperGraph(g1) && g1.IsSuperGraph(g)
}

func (g Graph) IsSuperGraph(g1 Graph) bool {
	return g.IsNodesSuperset(g1) && g.IsEdgesSuperset(g1) && g.IsConditionalEdgesSuperset(g1)
}

func (g Graph) IsNodesSuperset(g1 Graph) bool {
	for _, n := range g1.Nodes {
		if g.Find(n.Tag) == -1 {
			return false
		}
	}
	return true
}

func (g Graph) IsEdgesSuperset(g1 Graph) bool {
	for _, edge := range g1.Edges {
		i := g.Find(g1.Nodes[edge[0]].Tag)
		if i == -1 {
			return false
		}
		j := g.Find(g1.Nodes[edge[1]].Tag)
		if j == -1 {
			return false
		}
		if g.FindEdge(Edge{i, j}) == -1 {
			return false
		}
	}
	return true
}

func (g Graph) IsConditionalEdgesSuperset(g1 Graph) bool {
	for _, ce1 := range g1.ConditionalEdges {
		for _, e1 := range ce1.Edges {
			i := g.Find(e1.From)
			if i == -1 {
				return false
			}
			j := g.Find(e1.To)
			if j == -1 {
				return false
			}
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
}
