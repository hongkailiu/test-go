package cincinnati

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/util/sets"
)

func (g Graph) Equal(g1 Graph) bool {
	return g.IsSuperGraphOf(g1) && g1.IsSuperGraphOf(g)
}

func (g Graph) IsSuperGraphOf(g1 Graph, versions ...string) bool {
	return g.IsNodesSupersetOf(g1, versions...) && g.IsEdgesSupersetOf(g1, versions...) && g.IsConditionalEdgesSupersetOf(g1, versions...)
}

func (g Graph) IsNodesSupersetOf(g1 Graph, versions ...string) bool {
	if len(versions) > 0 {
		// It is possible: arch=amd64&channel=candidate-4.2&version=4.1.1
		if len(g.Nodes) == 0 && len(g1.Nodes) > 0 {
			logrus.WithField("versions", strings.Join(versions, ",")).WithField("g1.Nodes", len(g1.Nodes)).WithField("image", g1.Nodes[0].Image).Debug("Super graph has no node?")
		}
		return true
	}
	for i, n := range g1.Nodes {
		if g.FindNode(n) == -1 {
			logrus.WithField("nodeVersion", n.Version).WithField("nodeImage", n.Image).
				WithField("index", i).Error("node is not in graph")
			return false
		}
	}
	return true
}

func (g Graph) IsEdgesSupersetOf(g1 Graph, versions ...string) bool {
	versionSet := sets.New(versions...)
	for i, edge := range g1.Edges {
		if !versionSet.Has(g1.Nodes[edge[0]].Version.String()) {
			logrus.WithField("from", g1.Nodes[edge[0]].Version.String()).WithField("to", g1.Nodes[edge[1]].Version.String()).Debug("Ignored an irrelevant edge")
			continue
		}
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

func (n Node) Equals(n1 Node) bool {
	return n.Image == n1.Image && n.Version.Equals(n1.Version) &&
		equalChannels(n.Metadata[MetadataKeyChannels], n1.Metadata[MetadataKeyChannels]) &&
		n.Metadata[MetadataKeyArchitecture] == n1.Metadata[MetadataKeyArchitecture] &&
		n.Metadata["url"] == n1.Metadata["url"]
}

func equalChannels(s string, s1 string) bool {
	splits := strings.Split(s, ",")
	splits1 := strings.Split(s1, ",")
	return sets.New[string](splits...).Equal(sets.New[string](splits1...))
}

func (g Graph) FindNode(node Node) int {
	for i, n := range g.Nodes {
		if n.Equals(node) {
			return i
		}
	}
	return -1
}

func (g Graph) IsConditionalEdgesSupersetOf(g1 Graph, versions ...string) bool {
	versionSet := sets.New(versions...)
	for _, ce1 := range g1.ConditionalEdges {
		for _, e1 := range ce1.Edges {
			if !versionSet.Has(e1.From) {
				logrus.WithField("from", e1.From).WithField("to", e1.To).Debug("Ignored an irrelevant conditional edge")
				continue
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

	tests := []struct {
		name    string
		channel string
		version string
	}{
		{
			name:    "stable",
			channel: "stable-4.18",
			version: "4.18.10",
		},
		{
			name:    "candidate",
			channel: "candidate-4.18",
			version: "4.17.10",
		},
		{
			name:    "old candidate",
			channel: "candidate-4.2",
			version: "4.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arch := range []ArchParam{ArchParamAMD64, ArchParamARM64, ArchParamPPC64LE, ArchParamPPC64LE, ArchParamMULTI} {
				archStr := string(arch)
				verify(t, tt.channel, archStr, tt.version)
			}
		})
	}
}

func mustProduction(t *testing.T, channel, arch string) Graph {
	params := url.Values{}
	params.Add("channel", channel)
	params.Add("arch", arch)
	query := params.Encode()
	url, err := url.Parse("https://api.openshift.com/api/upgrades_info/graph")
	if err != nil {
		t.Fatalf("Failed to parse url: %v", err)
	}
	url.RawQuery = query

	production, err := getGraph(url.String())
	if err != nil {
		t.Fatalf("Failed to get graph: %v", err)
	}
	return production
}

func verify(t *testing.T, channel, arch string, versions ...string) {
	production := mustProduction(t, channel, arch)
	for _, version := range versions {
		// https://cincinnati-cincinnati-go.apps.ota-stage.q2z4.p1.openshiftapps.com/upgrades_info/v1/graph?channel=stable-4.10&arch=amd64&version=4.10.10
		params := url.Values{}
		params.Add("channel", channel)
		params.Add("arch", arch)
		params.Add("version", version)
		query := params.Encode()

		url, err := url.Parse("https://cincinnati-cincinnati-go.apps.ota-stage.q2z4.p1.openshiftapps.com/api/upgrades_info/graph")
		if err != nil {
			t.Fatalf("Failed to parse url: %v", err)
		}
		url.RawQuery = query

		graph, err := getGraph(url.String())
		if err != nil {
			t.Fatalf("Failed to get graph: %v", err)
		}

		if !production.IsSuperGraphOf(graph) {
			t.Error("production is not a super-graph")
		}
		if !graph.IsSuperGraphOf(production, version) {
			t.Error("cincinnati-go is not a super-graph")
		}
	}
}

func getVersions() (map[string][]string, error) {
	url := "https://raw.githubusercontent.com/openshift/cincinnati-graph-data/refs/heads/master/internal-channels/candidate.yaml"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var c Channel
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, err
	}

	result := map[string][]string{}
	for _, version := range c.Versions {
		v, err := semver.Parse(version)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%d.%d", v.Major, v.Minor)
		result[key] = append(result[key], version)
	}

	return result, nil

}

func TestIntegration_smoke(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") != "1" {
		t.Skip("integration tests skipped unless TEST_INTEGRATION=1")
	}

	tests := []struct {
		name        string
		channelType string
	}{
		{
			name:        "stable",
			channelType: "stable",
		},
		{
			name:        "fast",
			channelType: "fast",
		},
		{
			name:        "candidate",
			channelType: "candidate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions, err := getVersions()
			if err != nil {
				t.Fatalf("Failed to get versions: %v", err)
			}

			var keys []string
			for k := range versions {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				for _, arch := range []ArchParam{ArchParamAMD64, ArchParamARM64, ArchParamPPC64LE, ArchParamPPC64LE, ArchParamMULTI} {
					archStr := string(arch)
					verify(t, fmt.Sprintf("%s-%s", tt.name, k), archStr, versions[k]...)
				}
			}
		})
	}
}
