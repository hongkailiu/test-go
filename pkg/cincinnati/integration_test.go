package cincinnati

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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
			index := g.FindNode(n)
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
			logrus.WithField("index", i).WithField("edge", edge).
				WithField("g.from", from).WithField("g.to", to).
				WithField("from", g1.Nodes[edge[0]].Version.String()).
				WithField("to", g1.Nodes[edge[1]].Version.String()).
				Error("edge is not in graph")
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

func getGraph(url string) (Graph, error) {
	var g Graph
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return g, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return g, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return g, fmt.Errorf("error reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return g, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err := json.Unmarshal(raw, &g); err != nil {
		return g, err
	}

	return g, nil

}

func TestIntegration_dummy(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "1" {
		t.Skip("integration tests skipped unless TEST_INTEGRATION=1")
	}
	// https://cincinnati-cincinnati-go.apps.ota-stage.q2z4.p1.openshiftapps.com/upgrades_info/v1/graph?channel=stable-4.10&arch=amd64&version=4.10.10
	params := url.Values{}
	params.Add("channel", "stable-4.18")
	params.Add("arch", "amd64")
	params.Add("version", "4.18.10")
	query := params.Encode()

	url, err := url.Parse("https://cincinnati-cincinnati-go.apps.ota-stage.q2z4.p1.openshiftapps.com/upgrades_info/v1/graph")
	if err != nil {
		t.Fatalf("Failed to parse url: %v", err)
	}
	url.RawQuery = query

	graph, err := getGraph(url.String())
	if err != nil {
		t.Fatalf("Failed to get graph: %v", err)
	}

	url, err = url.Parse("https://api.openshift.com/api/upgrades_info/graph")
	if err != nil {
		t.Fatalf("Failed to parse url: %v", err)
	}
	url.RawQuery = query

	production, err := getGraph(url.String())
	if err != nil {
		t.Fatalf("Failed to get graph: %v", err)
	}

	if !production.IsSuperGraphOf(graph) {
		t.Fatal("production is not a super-graph of graph")
	}
}
