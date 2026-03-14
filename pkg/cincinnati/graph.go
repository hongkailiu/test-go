package cincinnati

import (
	"fmt"
	"reflect"

	"github.com/blang/semver/v4"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
)

type Graph struct {
	Nodes            []Node            `json:"nodes"`
	Edges            []Edge            `json:"edges"`
	ConditionalEdges []ConditionalEdge `json:"conditionalEdges"`

	Params GraphParams `json:"-"`
}

type Node struct {
	Version  semver.Version    `json:"version"`
	Image    string            `json:"payload"`
	Metadata map[string]string `json:"metadata,omitempty"`
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

type GraphHandlerFunc func(Graph) (Graph, error)

type GraphGenerator struct {
	handlers []GraphHandlerFunc
}

func generateEmptyGraph(params GraphParams) Graph {
	return Graph{
		Nodes:            []Node{},
		Edges:            []Edge{},
		ConditionalEdges: []ConditionalEdge{},
		Params:           params,
	}
}

func (g *GraphGenerator) Generate(params GraphParams) (Graph, error) {
	graph := generateEmptyGraph(params)
	var errs []error
	for _, h := range g.handlers {
		newGraph, err := h(graph)
		if err != nil {
			errs = append(errs, fmt.Errorf("error generating graph by %s: %w", reflect.TypeOf(h), err))
			return graph, kerrors.NewAggregate(errs)
		}
		graph = newGraph
	}
	return graph, kerrors.NewAggregate(errs)
}

func (g *GraphGenerator) Add(h GraphHandlerFunc) *GraphGenerator {
	g.handlers = append(g.handlers, h)
	return g
}

func NewGraphGenerator(handlers ...GraphHandlerFunc) GraphGenerator {
	return GraphGenerator{
		handlers: handlers,
	}
}
