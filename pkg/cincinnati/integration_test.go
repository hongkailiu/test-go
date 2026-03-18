package cincinnati

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TODO: compare result with production

func (g Graph) Equal(g1 Graph) bool {
	return g.IsSuperGraph(g1) && g1.IsSuperGraph(g)
}

func (g Graph) IsSuperGraph(g1 Graph) bool {
	return IsNodesSuperset(g.Nodes, g1.Nodes)
}

func (n Node) Equal(n1 Node) bool {
	return cmp.Equal(n, n1, cmpopts.IgnoreFields(Node{}, "Tag", "Previous"))
}

func IsNodesSuperset(small, big []Node) bool {
	for _, n := range small {
		if !n.isMemberOf(big) {
			return false
		}
	}
	return true
}

func (n Node) isMemberOf(nodes []Node) bool {
	for _, n1 := range nodes {
		if n1.Equal(n) {
			return true
		}
	}
	return false
}
