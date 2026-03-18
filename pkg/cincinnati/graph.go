package cincinnati

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/interrupts"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
)

// TODO: Generate Go struct from OpenAPI Specs or the other way around
// TODO: more test on arch and relevant suffix

type Graph struct {
	Nodes            []Node            `json:"nodes"`
	Edges            []Edge            `json:"edges"`
	ConditionalEdges []ConditionalEdge `json:"conditionalEdges"`

	Channels []string `json:"channels,omitempty"`
}

type Node struct {
	Version  semver.Version    `json:"version,omitempty"`
	Image    string            `json:"payload,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

	Tag      string   `json:"tag,omitempty"`
	Previous []string `json:"previous,omitempty"`
}

func (n Node) AddMetadata(k, v string) {
	if n.Metadata == nil {
		n.Metadata = map[string]string{}
	}
	n.Metadata[k] = v
}

func (n Node) DeleteMetadata(k string) {
	if _, ok := n.Metadata[k]; ok {
		delete(n.Metadata, k)
	}
}

type Edge [2]int

type ConditionalEdge struct {
	Edges []ConditionalUpdate     `json:"edges"`
	Risks []ConditionalUpdateRisk `json:"risks"`

	RisksKey string `json:"risksKey,omitempty"`
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
	var remove []int
	for i, node := range g.Nodes {
		if node.Version.LT(p.Version) {
			logrus.WithField("node.version", node.Version.String()).WithField("params.version", p.Version.String()).
				Debug("Ignored a smaller version")
			remove = append(remove, i)
			continue
		}
		if !archMatch(node.Tag, p.Arch) {
			logrus.WithField("node.version", node.Version.String()).WithField("node.tag", node.Tag).WithField("params.arch", p.Arch).
				Debug("Ignored a version not matching arch")
			remove = append(remove, i)
			continue
		}
		if channels := node.Metadata["io.openshift.upgrades.graph.release.channels"]; !strings.Contains(channels, p.Channel) {
			logrus.WithField("node.channels", node.Version.String()).WithField("params.channel", p.Channel).
				Debug("Ignored a version not in the channel")
			remove = append(remove, i)
			continue
		}
	}
	g = g.RemoveNodes(remove...)
	return p.directTargets(g), nil
}

func getTag(version, arch string) string {
	return fmt.Sprintf("%s%s", version, arch)
}

func archToTagSuffix(arch string) string {
	switch arch {
	case "amd64":
		return "-x86_64"
	case "arm64":
		return "-aarch64"
	default:
		return fmt.Sprintf("-%s", arch)
	}
}

func (p *GraphParams) directTargets(g Graph) Graph {
	keep := sets.New[int]()
	suffix := archToTagSuffix(p.Arch)
	fromTag := getTag(p.Version.String(), suffix)
	if i := g.Find(fromTag); i > -1 {
		logrus.WithField("tag", fromTag).WithField("i", i).Debug("Keep a node")
		keep.Insert(i)
	} else {
		logrus.WithField("version", p.Version.String()).WithField("arch", p.Arch).
			Debug("Could not find the node for the given params")
	}
	for _, edge := range g.Edges {
		if p.Version.Equals(g.Nodes[edge[0]].Version) {
			logrus.WithField("tag", g.Nodes[edge[1]].Tag).WithField("i", edge[1]).Debug("Keep a node")
			keep.Insert(edge[1])
		}
	}
	for _, edge := range g.ConditionalEdges {
		for _, riskEdge := range edge.Edges {
			if p.Version.String() == riskEdge.From {
				tag := getTag(riskEdge.To, suffix)
				if i := g.Find(tag); i > -1 {
					logrus.WithField("tag", tag).WithField("i", i).Debug("Keep a node")
					keep.Insert(i)
				}
			}
		}
	}
	var remove []int
	for i := range g.Nodes {
		if !keep.Has(i) {
			remove = append(remove, i)
		}
	}
	return g.RemoveNodes(remove...)
}

func (g Graph) RemoveNodes(remove ...int) Graph {
	if len(remove) == 0 {
		return g
	}
	logrus.WithField("nodes", len(g.Nodes)).
		WithField("edges", len(g.Edges)).
		WithField("conditionalEdges", len(g.ConditionalEdges)).
		WithField("remove", remove).
		Info("Removing nodes ...")
	removeVersions := sets.New[string]()
	for _, index := range remove {
		removeVersions.Insert(g.Nodes[index].Version.String())
	}
	var nodes []Node
	indexSet := sets.New[int](remove...)
	for i, node := range g.Nodes {
		if !indexSet.Has(i) {
			nodes = append(nodes, node)
		}
	}
	g.Nodes = nodes
	var conditionalEdges []ConditionalEdge
	for _, ce := range g.ConditionalEdges {
		var edges []ConditionalUpdate
		for _, edge := range ce.Edges {
			if !removeVersions.Has(edge.From) && !removeVersions.Has(edge.To) {
				edges = append(edges, edge)
			}
		}
		ce.Edges = edges
		if len(ce.Edges) > 0 {
			conditionalEdges = append(conditionalEdges, ce)
		}
	}
	g.Edges = g.RemoveEdges(remove)
	g.ConditionalEdges = conditionalEdges
	logrus.WithField("nodes", len(g.Nodes)).
		WithField("edges", len(g.Edges)).
		WithField("conditionalEdges", len(g.ConditionalEdges)).
		WithField("remove", remove).
		Info("Removed nodes")
	return g
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
	for _, s := range []string{"-x86_64", "-aarch64", "-s390x", "-ppc64le", "-multi"} {
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

func (g *GraphBuilder) Build(p GraphParams) (Graph, error) {
	zero := Graph{}.compatible()
	v, ok := g.cache.Get(cacheKeyOpenshiftUpgradeGraph)
	if !ok {
		return zero, fmt.Errorf("graph not found in cache")
	}
	graph := v.(Graph)
	var found bool
	for _, c := range graph.Channels {
		if c == p.Channel {
			found = true
		}
	}
	if !found {
		return zero, nil
	}
	shaped, err := p.shape(graph)
	if err != nil {
		return zero, err
	}
	return shaped.compatible(), nil
}

func (g Graph) compatible() Graph {
	for i := range g.Nodes {
		g.Nodes[i].Tag = ""
		g.Nodes[i].Previous = nil
	}
	if g.Nodes == nil {
		g.Nodes = []Node{}
	}
	if g.Edges == nil {
		g.Edges = []Edge{}
	}
	g.Channels = nil
	for i := range g.ConditionalEdges {
		g.ConditionalEdges[i].RisksKey = ""
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

		handles := []GraphHandlerFunc{g.repo.tagsToNodesAndEdges, gd.Shape}
		graph, err := buildOpenshiftUpgradeGraph(g.graphFile, handles, g.cache.Set, cache.NoExpiration)
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

func buildOpenshiftUpgradeGraph(graphFile string, handlers []GraphHandlerFunc, set func(k string, x interface{}, d time.Duration), d time.Duration) (Graph, error) {
	var graph Graph
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
				logrus.Info("Storing OpenShift upgrade graph loaded from file (to be refreshed if stale)")
				set(cacheKeyOpenshiftUpgradeGraph, graph, d)
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

func (g Graph) EnsureNode(n Node) Graph {
	if i := g.Find(n.Tag); i > -1 {
		g.Nodes[i] = n
	}
	g.Nodes = append(g.Nodes, n)
	return g
}

func (g Graph) EnsureEdges(edges []Edge) Graph {
	for _, edge := range edges {
		if g.FindEdge(edge) == -1 {
			g.Edges = append(g.Edges, edge)
		}
	}
	return g
}

func (g Graph) Find(tag string) int {
	for i, from := range g.Nodes {
		if tag == from.Tag {
			return i
		}
	}
	return -1
}

func (g Graph) FindEdge(edge Edge) int {
	for i, e := range g.Edges {
		if e == edge {
			return i
		}
	}
	return -1
}

func (g Graph) RemoveEdges(removedNodes []int) []Edge {
	removeSet := sets.New[int](removedNodes...)
	var edges []Edge
	for _, edge := range g.Edges {
		if removeSet.Has(edge[0]) || removeSet.Has(edge[1]) {
			continue
		}
		edges = append(edges, Edge{newIndex(removedNodes, edge[0]), newIndex(removedNodes, edge[1])})
	}
	return edges
}

func newIndex(removed []int, e int) int {
	if !sort.IntsAreSorted(removed) {
		panic(fmt.Sprintf("unsorted remove: %v", removed))
	}
	for i, r := range removed {
		if r > e {
			return e - i
		}
	}
	return e - len(removed)
}
