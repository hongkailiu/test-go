package cincinnati

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	Org  string
	Repo string
}

func (r *Repo) tags() ([]string, error) {
	// TODO https://oneuptime.com/blog/post/2026-02-08-how-to-list-all-tags-of-a-docker-image-on-docker-hub/view
	return []string{"4.10.0-rc.7-aarch64", "4.0.0-8", "4.0.0-0.okd-0", "4.19.3-x86_64"}, nil
}

func (r *Repo) tagsToNode(graph Graph) (Graph, error) {
	tags, err := r.tags()
	if err != nil {
		return Graph{}, fmt.Errorf("failed to fetch tags: %w", err)
	}
	for _, tag := range tags {
		v, err := semver.Parse(tag)
		if err != nil {
			logrus.WithField("tag", tag).Warnf("failed to parse tag %s (ignored): %v", tag, err)
			continue
		}
		if graph.Params.Version.GT(v) {
			logrus.WithField("version", v.String()).WithField("params.version", graph.Params.Version.String()).
				Debug("failed to parse tag %s (ignored): %v", tag, err)
			continue
		}
		graph.Nodes = append(graph.Nodes, Node{
			Version: v,
			Image:   "todo",
			Metadata: map[string]string{
				"io.openshift.upgrades.graph.release.channels":    "todo",
				"io.openshift.upgrades.graph.release.manifestref": "todo",
				"url": "todo",
			},
		})
	}
	return graph, nil
}

func (r *Repo) GetGraphHandlers() []GraphHandlerFunc {
	return []GraphHandlerFunc{r.tagsToNode}
}
