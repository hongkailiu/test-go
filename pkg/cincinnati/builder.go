package cincinnati

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/interrupts"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hongkailiu/test-go/pkg/util"
)

type GraphHandlerFunc func(context.Context, Graph) (Graph, error)

type GraphBuilder struct {
	graphFile    string
	graphDataDir string
	mockDir      string
	repo         *Repo

	cache Cache
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
			return false
		}
	}

	return true
}

func (g *GraphBuilder) Start(ctx context.Context) error {
	interrupts.TickLiteral(func() {
		logrus.Info("Building OpenShift upgrade graph ...")

		var gd CincinnatiGraphData

		if err := wait.PollUntilContextCancel(ctx, 3*time.Second, true, func(context.Context) (done bool, err error) {
			value, ok := g.cache.Get(cacheKeyCincinnatiGraphData)
			if !ok {
				logrus.Info("Loading Cincinnati graph data...")

				return false, nil
			}

			gd = value.(CincinnatiGraphData)

			return true, nil
		}); err != nil {
			logrus.WithError(err).Fatal("Failed to load graph data")
		}

		handles := []GraphHandlerFunc{g.repo.tagsToNodesAndEdges, gd.Shape, checkCycles}

		start := time.Now()
		graph, err := buildOpenshiftUpgradeGraph(ctx, g.graphFile, handles, g.cache.Set, cache.NoExpiration)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to build openshift upgrade graph")
		}
		d := time.Since(start)

		logrus.WithField("duration", d).
			WithField("nodes", len(graph.Nodes)).
			WithField("edges", len(graph.Edges)).
			WithField("conditionalEdges", len(graph.ConditionalEdges)).
			Info("Built OpenShift upgrade graph")
		g.cache.Set(cacheKeyOpenshiftUpgradeGraph, graph, cache.NoExpiration)
	}, 2*time.Hour)

	interrupts.TickLiteral(func() {
		err := g.writeOpenshiftUpgradeGraphToFile()
		if err != nil && !errors.Is(err, errGraphNotFoundInCache) {
			logrus.WithError(err).Error("Failed to write openshift upgrade graph to file")
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

func (g *GraphBuilder) writeOpenshiftUpgradeGraphToFile() error {
	v, ok := g.cache.Get(cacheKeyOpenshiftUpgradeGraph)
	if !ok {
		return errGraphNotFoundInCache
	}

	graph := v.(Graph)

	var raw []byte

	var err error

	if releaseMode() {
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

func buildOpenshiftUpgradeGraph(ctx context.Context, graphFile string, handlers []GraphHandlerFunc, set func(k string, x interface{}, d time.Duration), d time.Duration) (Graph, error) {
	var (
		graph  Graph
		loaded bool
	)

	if graphFile != "" {
		logrus.Info("Loading OpenShift upgrade graph from file ...")

		raw, err := util.ReadFileMaybeGZIP(graphFile)
		if err != nil {
			logrus.WithError(err).Warning("Failed to read openshift upgrade graph")
		} else {
			graphFromFile := Graph{}

			err = json.Unmarshal(raw, &graphFromFile)
			if err != nil {
				logrus.WithError(err).Warning("Failed to unmarshal openshift upgrade graph")
			} else {
				logrus.WithField("file", graphFile).Info("Loaded openshift upgrade graph from file")

				graph = graphFromFile
				loaded = true

				logrus.Info("Storing OpenShift upgrade graph loaded from file to cache (to be refreshed if stale)")
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
		newGraph, err := h(ctx, graph)
		if err != nil {
			errs = append(errs, fmt.Errorf("error generating graph by %s: %w", reflect.TypeOf(h), err))

			return graph, kerrors.NewAggregate(errs)
		}

		graph = newGraph
	}

	return graph, kerrors.NewAggregate(errs)
}

func NewGraphBuilder(file, graphDataDir, mockDir string, cache Cache, repo *Repo) *GraphBuilder {
	return &GraphBuilder{
		graphFile:    file,
		graphDataDir: graphDataDir,
		mockDir:      mockDir,
		cache:        cache,
		repo:         repo,
	}
}

// TODO: implement checkCycles
// checkCycles checks if the graph has a cycle
// It is enough to check the internal Field Previous because conditional update edges are derived from them.
func checkCycles(_ context.Context, g Graph) (Graph, error) {
	return g, nil
}
