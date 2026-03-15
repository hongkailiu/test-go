package cincinnati

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/interrupts"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
)

type Graph struct {
	Nodes            []Node            `json:"nodes"`
	Edges            []Edge            `json:"edges"`
	ConditionalEdges []ConditionalEdge `json:"conditionalEdges"`
}

type Node struct {
	Version  semver.Version    `json:"version,omitempty"`
	Image    string            `json:"payload,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

	Tag      string   `json:"tag,omitempty"`
	Previous []string `json:"previous,omitempty"`
}

func (g Graph) Equal(g1 Graph) bool {
	return g.IsSuperGraph(g1) && g1.IsSuperGraph(g)
}

func (g Graph) IsSuperGraph(g1 Graph) bool {
	// TODO
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

func (n Node) AddMetadata(k, v string) {
	if n.Metadata == nil {
		n.Metadata = map[string]string{}
	}
	n.Metadata[k] = v
}

type Edge [2]int

type ConditionalEdge struct {
	Edges []ConditionalUpdate     `json:"edges"`
	Risks []ConditionalUpdateRisk `json:"risks"`
}

type ConditionalUpdate struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ConditionalUpdateRisk struct {
	URL           string         `json:"url"`
	Name          string         `json:"name"`
	Message       string         `json:"message"`
	MatchingRules []MatchingRule `json:"matchingRules"`
}

type MatchingRule struct {
	Type   string       `json:"type"`
	PromQL *PromQLQuery `json:"promql,omitempty"`
}

type PromQLQuery struct {
	PromQL string `json:"promql"`
}

type ChannelInfo struct {
	Name        string
	Description string
	Example     string
	CurlCommand string
}

type GraphParams struct {
	Arch    string
	Channel string
	Id      string
	Version semver.Version
}

func (p *GraphParams) shape(g Graph) (Graph, error) {
	var nodes []Node
	for _, node := range g.Nodes {
		if node.Version.LT(p.Version) {
			logrus.WithField("node.version", node.Version.String()).WithField("params.version", p.Version.String()).
				Debug("Ignored a smaller version")
			continue
		}
		if !archMatch(node.Tag, p.Arch) {
			continue
		}
		node.Tag = ""
		node.Previous = nil
		nodes = append(nodes, node)
	}
	g.Nodes = nodes
	// TODO: use p to shape graph
	return g, nil
}

func archMatch(tag string, arch string) bool {
	switch arch {
	case "amd64":
		return strings.HasSuffix(tag, "-x86_64") || !hasArchSuffix(tag)
	case "arm64":
		return strings.HasSuffix(tag, "-aarch64")
	default:
		return strings.HasSuffix(tag, fmt.Sprintf("-%s", arch))
	}
}

func hasArchSuffix(tag string) bool {
	for _, s := range []string{"-x86_64", "-aarch64", "-s390x", "-ppc64le"} {
		if strings.HasSuffix(tag, s) {
			return true
		}
	}
	return false
}

type GraphHandlerFunc func(Graph) (Graph, error)

type GraphBuilder struct {
	ctx          context.Context
	graphFile    string
	graphDataDir string
	mockDir      string
	repo         *Repo

	cache Cache
}

func generateEmptyGraph() Graph {
	return Graph{
		Nodes:            []Node{},
		Edges:            []Edge{},
		ConditionalEdges: []ConditionalEdge{},
	}
}

func (g *GraphBuilder) Build(p GraphParams) (Graph, error) {
	var zero Graph
	v, ok := g.cache.Get(cacheKeyOpenshiftUpgradeGraph)
	if !ok {
		return zero, fmt.Errorf("graph not found in cache")
	}
	shaped, err := p.shape(v.(Graph))
	if err != nil {
		return zero.compatible(), err
	}
	return shaped.compatible(), nil
}

func (g Graph) compatible() Graph {
	if g.Nodes == nil {
		g.Nodes = []Node{}
	}
	if g.Edges == nil {
		g.Edges = []Edge{}
	}
	if g.ConditionalEdges == nil {
		g.ConditionalEdges = []ConditionalEdge{}
	}
	return g
}

const cacheKeyOpenshiftUpgradeGraph = "openshift-upgrade-graph"
const cacheKeyCincinnatiGraphData = "cincinnati-graph-data"

func (g *GraphBuilder) Ready() bool {
	for _, key := range []string{cacheKeyOpenshiftUpgradeGraph, cacheKeyCincinnatiGraphData} {
		_, ok := g.cache.Get(key)
		if !ok {
			return false
		}
	}
	return true
}

func (g *GraphBuilder) Start() error {
	interrupts.TickLiteral(func() {
		logrus.Info("Building OpenShift upgrade graph ...")

		var gd CincinnatiGraphData
		if err := wait.PollUntilContextCancel(g.ctx, 3*time.Second, true, func(context.Context) (done bool, err error) {
			value, ok := g.cache.Get(cacheKeyCincinnatiGraphData)
			if !ok {
				logrus.Info("Loading Cincinnati graph data...")
				return false, nil
			}
			gd = value.(CincinnatiGraphData)
			return true, nil
		}); err != nil {
			logrus.WithError(err).Warn("Cincinnati graph data not loaded (using empty instead)")
		}

		handles := []GraphHandlerFunc{g.repo.tagsToNodes, gd.Shape}
		graph, err := buildOpenshiftUpgradeGraph(g.graphFile, handles)
		if err != nil {
			logrus.WithError(err).Error("Failed to build openshift upgrade graph")
		}
		logrus.Info("Built OpenShift upgrade graph")
		g.cache.Set(cacheKeyOpenshiftUpgradeGraph, graph, cache.NoExpiration)
	}, 2*time.Hour)

	interrupts.TickLiteral(func() {
		if err := g.storeOpenshiftUpgradeGraph(); err != nil {
			logrus.WithError(err).Error("Failed to store openshift upgrade graph")
		}
	}, time.Minute)

	interrupts.TickLiteral(func() {
		dir := g.graphDataDir
		if g.mockDir != "" {
			dir = filepath.Join(g.mockDir, "graph-data")
		}
		gd, err := LoadGraphData(dir)
		if err != nil {
			logrus.WithError(err).WithField("dir", dir).Error("Failed to load graph data")
			return
		}
		g.cache.Set(cacheKeyCincinnatiGraphData, *gd, 10*time.Minute)
	}, time.Minute)

	return nil
}

func (g *GraphBuilder) storeOpenshiftUpgradeGraph() error {
	v, ok := g.cache.Get(cacheKeyOpenshiftUpgradeGraph)
	if !ok {
		return fmt.Errorf("graph not found in cache")
	}
	graph := v.(Graph)
	raw, err := json.Marshal(graph)
	if err != nil {
		return fmt.Errorf("error serializing graph: %w", err)
	}
	err = os.WriteFile(g.graphFile, raw, 0644)
	if err != nil {
		return fmt.Errorf("error write graph to file %s: %w", g.graphFile, err)
	}
	return nil
}

func buildOpenshiftUpgradeGraph(graphFile string, handlers []GraphHandlerFunc) (Graph, error) {
	graph := generateEmptyGraph()
	var loaded bool
	if graphFile != "" {
		logrus.Info("Loading OpenShift upgrade graph from file ...")
		raw, err := os.ReadFile(graphFile)
		if err != nil {
			logrus.WithError(err).Warning("Failed to read openshift upgrade graph")
		} else {
			graphFromFile := Graph{}
			if err = json.Unmarshal(raw, &graphFromFile); err != nil {
				logrus.WithError(err).Warning("Failed to unmarshal openshift upgrade graph")
			} else {
				logrus.WithField("file", graphFile).Info("Loaded openshift upgrade graph from file")
				graph = graphFromFile
				loaded = true
			}
		}
	}

	if !loaded {
		logrus.Info("Building OpenShift upgrade graph from scratch ...")
	} else {
		logrus.Info("Building OpenShift upgrade graph incrementally ...")
	}

	var errs []error
	for _, h := range handlers {
		newGraph, err := h(graph)
		if err != nil {
			errs = append(errs, fmt.Errorf("error generating graph by %s: %w", reflect.TypeOf(h), err))
			return graph, kerrors.NewAggregate(errs)
		}
		graph = newGraph
	}
	return graph, kerrors.NewAggregate(errs)
}

func NewGraphBuilder(ctx context.Context, file, graphDataDir, mockDir string, cache Cache, repo *Repo) *GraphBuilder {
	return &GraphBuilder{
		ctx:          ctx,
		graphFile:    file,
		graphDataDir: graphDataDir,
		mockDir:      mockDir,
		cache:        cache,
		repo:         repo,
	}
}

func (n Node) getFrom(tag string, p string, graph Graph) int {
	arch := getArch(tag)
	for i, node := range graph.Nodes {
		if node.Version.String() == p {
			if arch1 := getArch(node.Tag); arch1 == arch || arch1 == "" {
				return i
			}
		}
	}
	return -1
}

func getArch(tag string) string {
	if strings.HasSuffix(tag, "-x86_64") {
		return "amd64"
	} else if strings.HasSuffix(tag, "-aarch64") {
		return "arm64"
	} else if strings.HasSuffix(tag, "-s390x") {
		return "s390x"
	} else if strings.HasSuffix(tag, "-ppc64le") {
		return "ppc64le"
	} else if strings.HasSuffix(tag, "-multi") {
		return "multi"
	}
	return ""
}

func (n Node) getPrevious(graph Graph) []Edge {
	i := graph.Find(n.Tag)
	if i == -1 {
		logrus.WithField("tag", n.Tag).Warning("Found no node with the tag in the graph")
		return nil
	}
	to := i

	var edges []Edge
	for _, p := range n.Previous {
		i := n.getFrom(n.Tag, p, graph)
		if i == -1 {
			continue
		}
		from := i
		edges = append(edges, Edge{from, to})
	}
	return edges
}

func (g Graph) Find(tag string) int {
	for i, from := range g.Nodes {
		if tag == from.Tag {
			return i
		}
	}
	return -1
}

func EnsureEdges(g *Graph, edges []Edge) {
	if g == nil {
		panic("nil graph cannot not contain any edges")
	}
	for _, edge := range edges {
		if g.FindEdge(edge) == -1 {
			g.Edges = append(g.Edges, edge)
		}
	}
}

func (g Graph) FindEdge(edge Edge) int {
	for i, e := range g.Edges {
		if e == edge {
			return i
		}
	}
	return -1
}
