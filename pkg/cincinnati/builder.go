package cincinnati

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hongkailiu/test-go/pkg/util"
)

type GraphHandlerFunc func(context.Context, Graph) (Graph, error)

type GraphHandler struct {
	HandlerFunc GraphHandlerFunc
	Name        string
}

type GraphBuilder struct {
	graphFile    string
	graphDataDir string
	repo         *Repo

	cache Cache

	releaseMode bool
}

var errGraphNotFoundInCache = errors.New("graph not found in cache")

func (g *GraphBuilder) Build(p GraphParams) (Graph, error) {
	zero := Graph{}.compatible()

	v, ok := g.cache.Get(cacheKeyOpenshiftUpgradeGraph)
	if !ok {
		return zero, errGraphNotFoundInCache
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

	// multi-arch is GA on 4.13 (TP within 4.11&4.12) but Cincinnati shows upgrade paths from 4.12+
	if p.Arch == Multi && p.Version.Major == 4 && p.Version.Minor <= 11 {
		return zero, nil
	}

	shaped, err := p.shape(graph)
	if err != nil {
		return zero, err
	}

	return shaped.compatible(), nil
}

const cacheKeyOpenshiftUpgradeGraph = "openshift-upgrade-graph"
const cacheKeyCincinnatiGraphData = "cincinnati-graph-data"

func (g *GraphBuilder) Ready() bool {
	for _, key := range []string{cacheKeyOpenshiftUpgradeGraph, cacheKeyCincinnatiGraphData} {
		_, ok := g.cache.Get(key)
		if !ok {
			logrus.WithField("missing", key).Info("Sever is not ready yet")
			return false
		}
	}

	return true
}

func (g *GraphBuilder) CacheGraph(ctx context.Context) error {
	logrus.Info("Building OpenShift upgrade graph ...")

	var gd CincinnatiGraphData

	if err := wait.PollUntilContextCancel(ctx, 3*time.Second, true, func(context.Context) (done bool, err error) {
		value, ok := g.cache.Get(cacheKeyCincinnatiGraphData)
		if !ok {
			logrus.Info("Waiting for loading Cincinnati graph data...")

			return false, nil
		}

		gd = value.(CincinnatiGraphData)

		return true, nil
	}); err != nil {
		logrus.WithError(err).Fatal("Failed to load graph data")
	}

	handles := []GraphHandler{
		{
			Name:        "repo.tagsToNodesAndEdges",
			HandlerFunc: g.repo.tagsToNodesAndEdges,
		},
		{
			Name:        "graphData.shape",
			HandlerFunc: gd.Shape,
		},
	}

	start := time.Now()
	graph, err := buildOpenshiftUpgradeGraph(ctx, g.graphFile, handles, func(graph Graph) {
		g.cache.Set(cacheKeyOpenshiftUpgradeGraph, graph, cache.NoExpiration)
		if err := g.writeOpenshiftUpgradeGraphToFile(graph); err != nil {
			logrus.WithError(err).Error("Failed to write openshift upgrade graph to file")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to build OpenShift upgrade graph: %w", err)
	}
	d := time.Since(start)

	logrus.WithField("duration", d).
		WithField("nodes", len(graph.Nodes)).
		WithField("edges", len(graph.Edges)).
		WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Built OpenShift upgrade graph")

	return nil
}

func (g *GraphBuilder) CacheGraphData() error {
	logrus.Info("Loading graph data ...")
	dir := g.graphDataDir

	gd, err := LoadGraphData(dir)
	if err != nil {
		logrus.WithError(err).WithField("dir", dir).Error("Failed to load graph data")
		return fmt.Errorf("failed to load graph data: %w", err)
	}

	g.cache.Set(cacheKeyCincinnatiGraphData, *gd, 10*time.Minute)
	logrus.WithField("blocked", len(gd.BlockedEdges)).
		WithField("channels", len(gd.Channels)).
		WithField("removed", len(gd.RemovedEdges)).Info("Loaded graph data")
	return nil
}

func (g *GraphBuilder) writeOpenshiftUpgradeGraphToFile(graph Graph) error {

	var raw []byte
	var err error

	if g.releaseMode {
		raw, err = json.Marshal(graph)
	} else {
		raw, err = json.MarshalIndent(graph, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("error serializing graph: %w", err)
	}

	err = util.WriteBytesMaybeGZIP(g.graphFile, raw)
	if err != nil {
		return fmt.Errorf("error write graph to file %s: %w", g.graphFile, err)
	}

	return nil
}

func buildOpenshiftUpgradeGraph(ctx context.Context, graphFile string, handlers []GraphHandler, cacheGraph func(graph Graph)) (Graph, error) {
	var (
		graph  Graph
		loaded bool
	)

	if graphFile != "" {
		logrus.WithField("graphFile", graphFile).Info("Loading OpenShift upgrade graph from file ...")

		raw, err := util.ReadFileMaybeGZIP(graphFile)
		if err != nil {
			logrus.WithError(err).Warning("Failed to read openshift upgrade graph")
		} else {
			graphFromFile := Graph{}

			err = json.Unmarshal(raw, &graphFromFile)
			if err != nil {
				logrus.WithError(err).Warning("Failed to unmarshal openshift upgrade graph")
			} else {
				logrus.WithField("graphFile", graphFile).Info("Loaded openshift upgrade graph from file")

				graph = graphFromFile
				loaded = true
				logrus.Info("Storing OpenShift upgrade graph loaded from file to cache (to be refreshed if stale)")
				cacheGraph(graph)
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
		start := time.Now()
		logrus.WithField("handler", h.Name).Info("handler started")
		newGraph, err := h.HandlerFunc(ctx, graph)
		if err != nil {
			errs = append(errs, fmt.Errorf("error generating graph by %s: %w", h.Name, err))

			return graph, kerrors.NewAggregate(errs)
		}

		graph = newGraph
		d := time.Since(start)
		logrus.WithField("handler", h.Name).WithField("duration", d).Info("handler completed")

	}

	cacheGraph(graph)
	return graph, kerrors.NewAggregate(errs)
}

func NewGraphBuilder(file, graphDataDir string, cache Cache, repo *Repo, releaseMode bool) *GraphBuilder {
	logrus.WithField("graphFile", file).
		WithField("graphDataDir", graphDataDir).
		Info("Create Graph Builder")
	return &GraphBuilder{
		graphFile:    file,
		graphDataDir: graphDataDir,
		cache:        cache,
		repo:         repo,
		releaseMode:  releaseMode,
	}
}
