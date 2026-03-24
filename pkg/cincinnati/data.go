package cincinnati

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type BlockedEdge struct {
	RemovedEdge

	FixedIn       string         `json:"fixedIn,omitempty"`
	URL           string         `json:"url,omitempty"`
	Name          string         `json:"name,omitempty"`
	Message       string         `json:"message,omitempty"`
	MatchingRules []MatchingRule `json:"matchingRules,omitempty"`
}

type Feeder struct {
	Name   string
	Errata string
	Filter string
}

type Channel struct {
	Name     string
	Versions []string
}

// RemovedEdge is an edge removed from the upgrade graph.
type RemovedEdge struct {
	To   string
	From string

	FromRegex *regexp.Regexp `json:"-"`
}
type CincinnatiGraphData struct {
	BlockedEdges []BlockedEdge `json:"blockedEdges,omitempty"`

	Channels     []Channel     `json:"channels,omitempty"`
	RemovedEdges []RemovedEdge `json:"removedEdges,omitempty"`
}

func LoadGraphData(dir string) (*CincinnatiGraphData, error) {
	var graphData CincinnatiGraphData

	channelsDir := filepath.Join(dir, "channels")

	err := filepath.WalkDir(channelsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		raw, err := os.ReadFile(filepath.Join(channelsDir, d.Name()))
		if err != nil {
			return err
		}

		var c Channel
		if err := yaml.Unmarshal(raw, &c); err != nil {
			return err
		}

		graphData.Channels = append(graphData.Channels, c)

		return nil
	})
	if err != nil {
		return nil, err
	}

	blockedEdgesDir := filepath.Join(dir, "blocked-edges")

	err = filepath.WalkDir(blockedEdgesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		raw, err := os.ReadFile(filepath.Join(blockedEdgesDir, d.Name()))
		if err != nil {
			return err
		}

		var be BlockedEdge
		if err := yaml.Unmarshal(raw, &be); err != nil {
			return err
		}

		if len(be.MatchingRules) == 0 {
			graphData.RemovedEdges = append(graphData.RemovedEdges, RemovedEdge{From: be.From, To: be.To})

			return nil
		}

		graphData.BlockedEdges = append(graphData.BlockedEdges, be)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &graphData, nil
}

func (gd CincinnatiGraphData) Shape(_ context.Context, graph Graph) (Graph, error) {
	var remove []int

	for i, node := range graph.Nodes {
		channels := gd.listChannels(node.Version.String())
		if len(channels) > 0 {
			graph.Channels = sets.List[string](sets.New[string](graph.Channels...).Insert(channels...))
			graph.Nodes[i].SetMetadata(MetadataKeyChannels, strings.Join(channels, ","))
		} else {
			delete(graph.Nodes[i].Metadata, MetadataKeyChannels)
			logrus.WithField("index", i).WithField("tag", node.Tag).WithField("version", node.Version.String()).Debug("Node in no channels to remove")
			remove = append(remove, i)
		}
	}

	if len(remove) > 0 {
		graph = graph.RemoveNodes(remove...)
	}

	var edges []Edge

	// It is ensured that each conditional risk edge (from, to)
	// occurs only once among all conditional updates
	graph.ConditionalEdges = nil
	for _, edge := range graph.Edges {
		if removed := gd.IsRemovedEdge(graph.Nodes[edge[0]], graph.Nodes[edge[1]].Version.String()); removed {
			continue
		}
		if ce := gd.BecomeConditional(graph.Nodes[edge[0]], graph.Nodes[edge[1]].Version.String()); len(ce.Risks) > 0 {
			var found bool

			for i, exiting := range graph.ConditionalEdges {
				if ce.RisksKey == exiting.RisksKey {
					found = true
					graph.ConditionalEdges[i].Edges = ensureConditionalEdges(graph.ConditionalEdges[i].Edges, ce.Edges)
				}
			}

			if !found {
				graph.ConditionalEdges = append(graph.ConditionalEdges, ce)
			}

			continue
		}

		edges = append(edges, edge)
	}

	graph.Edges = edges

	return graph, nil
}

func ensureConditionalEdges(existing []ConditionalUpdate, new []ConditionalUpdate) []ConditionalUpdate {
	for _, edge := range new {
		found := false
		for _, existingEdge := range existing {
			if existingEdge.From == edge.From && existingEdge.To == edge.To {
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, edge)
		}
	}
	return existing
}

func (gd CincinnatiGraphData) listChannels(version string) []string {
	var channels []string

	for _, channel := range gd.Channels {
		for _, v := range channel.Versions {
			if v == version {
				channels = append(channels, channel.Name)
			}
		}
	}

	return channels
}

func (gd CincinnatiGraphData) BecomeConditional(from Node, to string) ConditionalEdge {
	var (
		risks []ConditionalUpdateRisk
		key   []int
	)

	for i, blockedEdge := range gd.BlockedEdges {
		if to != blockedEdge.To {
			continue
		}

		if blockedEdge.FromRegex == nil {
			fromRegex, err := regexp.Compile(blockedEdge.From)
			if err != nil {
				logrus.WithError(err).
					WithField("name", blockedEdge.Name).
					WithField("from", blockedEdge.From).
					WithField("to", blockedEdge.To).
					Error("Failed to compile regex")

				continue
			}

			blockedEdge.FromRegex = fromRegex
		}

		if blockedEdge.FromRegex.MatchString(from.Tag) {
			logrus.WithField("from", from.Tag).WithField("to", to).Debug("Found blocked edge")

			key = append(key, i)
			risks = append(risks, ConditionalUpdateRisk{
				Name:          blockedEdge.Name,
				URL:           blockedEdge.URL,
				Message:       blockedEdge.Message,
				MatchingRules: blockedEdge.MatchingRules,
			})
		}
	}

	return ConditionalEdge{
		Risks: risks,
		Edges: []ConditionalUpdate{
			{
				From: from.Version.String(),
				To:   to,
			},
		},
		RisksKey: fmt.Sprintf("%v", key),
	}
}

func (gd CincinnatiGraphData) IsRemovedEdge(from Node, to string) bool {
	for _, removedEdge := range gd.RemovedEdges {
		if to != removedEdge.To {
			continue
		}

		if removedEdge.FromRegex == nil {
			fromRegex, err := regexp.Compile(removedEdge.From)
			if err != nil {
				logrus.WithError(err).
					WithField("from", removedEdge.From).
					WithField("to", removedEdge.To).
					Error("Failed to compile regex")

				continue
			}

			removedEdge.FromRegex = fromRegex
		}

		if removedEdge.FromRegex.MatchString(from.Tag) {
			logrus.WithField("from", from.Tag).WithField("to", to).Debug("Found removed edge")
			return true
		}
	}
	return false
}
