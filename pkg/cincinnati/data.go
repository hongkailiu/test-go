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
	To            string
	From          string
	FixedIn       string
	URL           string
	Name          string
	Message       string
	MatchingRules []MatchingRule
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

type CincinnatiGraphData struct {
	BlockedEdges []BlockedEdge
	Channels     []Channel
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
			graph.Nodes[i].AddMetadata(MetadataKeyChannels, strings.Join(channels, ","))
		} else {
			delete(graph.Nodes[i].Metadata, MetadataKeyChannels)
			remove = append(remove, i)
		}
	}

	if len(remove) > 0 {
		graph = graph.RemoveNodes(remove...)
	}

	var edges []Edge

	graph.ConditionalEdges = nil
	for _, edge := range graph.Edges {
		if ce := gd.BecomeConditional(graph.Nodes[edge[0]], graph.Nodes[edge[1]].Version.String()); len(ce.Risks) > 0 {
			var found bool

			for i, exiting := range graph.ConditionalEdges {
				if ce.RisksKey == exiting.RisksKey {
					found = true

					graph.ConditionalEdges[i].Edges = append(graph.ConditionalEdges[i].Edges, ce.Edges...)
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

		re, err := regexp.Compile(blockedEdge.From)
		if err != nil {
			logrus.WithError(err).
				WithField("name", blockedEdge.Name).
				WithField("from", blockedEdge.From).
				WithField("to", blockedEdge.To).
				Error("Failed to compile regex")

			continue
		}

		if re.MatchString(from.Tag) {
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
