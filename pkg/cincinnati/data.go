package cincinnati

import (
	"io/fs"
	"os"
	"path/filepath"

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
	Feeder     Feeder
	Name       string
	Versions   []string
	Tombstones []string
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

func (gd CincinnatiGraphData) Shape(graph Graph) (Graph, error) {
	// TODO
	return graph, nil
}
