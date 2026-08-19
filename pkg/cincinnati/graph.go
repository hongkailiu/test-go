package cincinnati

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/sets"
)

// TODO: Generate Go struct from OpenAPI Specs or the other way around

// multi-arch since 4.3.14: tag 4.2.1 -> 4.3.14-x86_64 with version 4.3.14
// condition update since 4.7?: default blocking -> conditional blocking with MatchRules
// tag could have arch suffix before 4.3 before 4.3, 4.2.11-s390x with version 4.2.11-s390x

// Cincinnati supports OKD as well.
// The release repo is quay.io/okd/scos-release.
// - The first tag is 4.12.0-0.okd-scos-2022-10-22-232744
// - The first manifest tag is 4.16.0-0.okd-scos-2024-07-29-144356 which has only 1 shard for amd64
// - The first manifest tag with more than 1 shard are 4.22.0-okd-scos.6 and 5.0.0-okd-scos.ec.4: amd64 and arm64
// - The first upgrade path is from 4.19.0-okd-scos.19 to 4.20.0-okd-scos.0
// OKD has its own graph-data
// https://github.com/okd-project/cincinnati-graph-data
// - It has only stable channels and starts from stable-4.20

type Graph struct {
	Version          int               `json:"version"`
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
	Arch     string   `json:"arch,omitempty"`
}

func SetMetadata(n *Node, k, v string) {
	if n.Metadata == nil {
		n.Metadata = map[string]string{}
	}

	n.Metadata[k] = v
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

		if node.Arch != p.Arch {
			logrus.WithField("node.version", node.Version.String()).WithField("node.tag", node.Tag).WithField("params.arch", p.Arch).
				Debug("Ignored a version not matching arch")

			remove = append(remove, i)

			continue
		}

		if channels := node.Metadata[MetadataKeyChannels]; !strings.Contains(channels, p.Channel) {
			logrus.WithField("node.channels", node.Version.String()).WithField("params.channel", p.Channel).
				Debug("Ignored a version not in the channel")

			remove = append(remove, i)

			continue
		}
	}

	g = g.RemoveNodes(remove...)

	return p.directTargets(g), nil
}

func (g Graph) getTagAndIndex(version string, arch string) (string, int) {

	for _, node := range g.Nodes {
		if node.Version.String() == version && node.Arch == arch {
			if i := g.Find(node.Tag); i > -1 {
				return node.Tag, i
			}
		}
	}
	return "", -1
}

func (p *GraphParams) directTargets(g Graph) Graph {
	keep := sets.New[int]()

	fromTag, i := g.getTagAndIndex(p.Version.String(), p.Arch)
	if i > -1 {
		logrus.WithField("tag", fromTag).WithField("i", i).Debug("Keep a node")
		keep.Insert(i)
	} else {
		// The client does the same.
		// https://github.com/openshift/cluster-version-operator/blob/55fa1518aa72b0b243584bf9a07c73810617f261/pkg/cincinnati/cincinnati.go#L235
		logrus.WithField("version", p.Version.String()).WithField("arch", p.Arch).
			Debug("Could not find the node for the given params and thus returned the zero graph")
		return Graph{}
	}

	for _, edge := range g.Edges {
		if p.Version.Equals(g.Nodes[edge[0]].Version) {
			logrus.WithField("tag", g.Nodes[edge[1]].Tag).WithField("i", edge[1]).Debug("Keep a node for edge")
			keep.Insert(edge[1])
		} else {
			logrus.WithField("tag", g.Nodes[edge[1]].Tag).WithField("i", edge[1]).Debug("Remove a node for edge")
		}
	}

	for _, edge := range g.ConditionalEdges {
		for _, riskEdge := range edge.Edges {
			if p.Version.String() == riskEdge.From {
				tag, i := g.getTagAndIndex(riskEdge.To, p.Arch)
				if i > -1 {
					logrus.WithField("tag", tag).WithField("i", i).Debug("Keep a node for conditional edge")
					keep.Insert(i)
				} else {
					logrus.WithField("tag", tag).WithField("i", i).Debug("remove a node for conditional edge")
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
		WithField("remove", len(remove)).
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

	// ConditionalEdge is not arch-specific and thus cannot be deleted unless no nodes with any arch use it.
	for _, node := range g.Nodes {
		if version := node.Version.String(); removeVersions.Has(version) {
			removeVersions.Delete(version)
		}
	}

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
		WithField("remove", len(remove)).
		Info("Removed nodes")

	return g
}

const Multi = "multi"

func (g Graph) compatible() Graph {
	for i := range g.Nodes {
		g.Nodes[i].Tag = ""
		g.Nodes[i].Previous = nil
		g.Nodes[i].Arch = ""
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
	if g.Version == 0 {
		g.Version = 1
	}

	return g
}

// getFrom return the index of the node which is with the given version and the same arch with the node in the given graph, or
// -1 if such a node cannot be found.
func (n Node) getFrom(version string, graph Graph) int {
	for i, node := range graph.Nodes {
		versionStr := node.Version.String()
		if versionStr == version {
			if node.Arch == n.Arch {
				return i
			}
		}
	}

	return -1
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
		i := n.getFrom(p, graph)
		if i == -1 {
			continue
		}

		from := i
		edges = append(edges, Edge{from, to})
	}

	return edges
}

func (g Graph) EnsureNode(n Node) Graph {
	if err := n.validPrevious(); err != nil {
		logrus.WithField("tag", n.Tag).WithError(err).Error("Ignored an invalid previous node")
		return g
	}
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

func (n Node) validPrevious() error {
	// Invalid previous should never happen, but we check it anyway for robustness.
	// This is a stronger checking than graph having no cycles.
	for i, p := range n.Previous {
		pVersion, err := semver.Parse(p)
		if err != nil {
			return fmt.Errorf("error parsing previous version %s of tag %s: %w", p, n.Tag, err)
		}
		if !pVersion.LT(n.Version) {
			logrus.WithField("index", i).WithField("p", p).
				WithField("tag", n.Tag).WithField("version", n.Version.String()).
				Error("previous is not smaller")
			return fmt.Errorf("tag's previous are not always smaller: %s", n.Tag)
		}
	}
	return nil
}
