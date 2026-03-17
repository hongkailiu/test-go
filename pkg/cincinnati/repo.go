package cincinnati

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type Repo struct {
	ctx            context.Context
	client         Client
	registry       string
	repo           string
	mockDir        string
	cache          Cache
	queue          Queue
	maxConcurrency int
}

func NewRepo(ctx context.Context, client Client, registry, repo, mockDir string, maxConcurrency int, cache Cache) *Repo {
	return &Repo{
		ctx:            ctx,
		client:         client,
		registry:       registry,
		repo:           repo,
		mockDir:        mockDir,
		cache:          cache,
		maxConcurrency: maxConcurrency,
	}
}

type Queue interface {
	Add(string)
	Get() (string, bool)
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

type TagsListData struct {
	Tags []string `json:"tags"`
}

func (r *Repo) tags() ([]string, error) {
	key := fmt.Sprintf("tagsWithCache-%s/%s", r.registry, r.repo)
	if v, ok := r.cache.Get(key); ok {
		return v.([]string), nil
	}
	if r.mockDir != "" {
		raw, err := os.ReadFile(filepath.Join(r.mockDir, "repo.tags.list.json"))
		if err != nil {
			return nil, err
		}
		data := TagsListData{}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		r.cache.Set(fmt.Sprintf("tagsWithCache-%s/%s", r.registry, r.repo), data.Tags, time.Duration(0))
		return data.Tags, nil
	}

	var ret []string
	url := fmt.Sprintf("%s/v2/%s/tags/list", r.registry, r.repo)
	for {
		tags, next, err := allTags(r.client, url)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tags...)
		if next != "" {
			nextURL, err := getNextURL(next)
			if err != nil {
				return nil, err
			}
			url = fmt.Sprintf("%s%s", r.registry, nextURL)
		} else {
			break
		}
	}
	logrus.WithField("key", key).WithField("value", ret).Info("Cache set")
	r.cache.Set(key, ret, time.Duration(0))
	return ret, nil
}

func getNextURL(next string) (string, error) {
	// </v2/openshift-release-dev/ocp-release/tags/list?n=100&last=4.10.36-aarch64>; rel="next"
	splits := strings.Split(next, ";")
	if len(splits) != 2 {
		return "", fmt.Errorf("invalid next: %s", next)
	}
	return strings.Trim(splits[0], "<>"), nil
}

func allTags(client Client, url string) ([]string, string, error) {
	// Reference https://oneuptime.com/blog/post/2026-02-08-how-to-list-all-tags-of-a-docker-image-on-docker-hub/view
	// https://quay.io/v2/openshift-release-dev/ocp-release/tags/list
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code %d for url %s", res.StatusCode, req.URL)
	}

	link := res.Header.Get("link")
	if link != "" {
		logrus.WithField("link", link).Debugf("Found link")
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	data := TagsListData{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, "", err
	}
	return data.Tags, link, nil
}

func (r *Repo) tagsToNodesAndEdges(graph Graph) (Graph, error) {
	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraping the repository for nodes and edges ...")
	tags, err := r.tags()
	if err != nil {
		return Graph{}, fmt.Errorf("failed to fetch tags: %w", err)
	}
	sem := semaphore.NewWeighted(int64(r.maxConcurrency))
	var missed []string
	for _, tag := range tags {
		var found bool
		// TODO: define the func on the struct
		for _, node := range graph.Nodes {
			if tag == node.Tag {
				found = true
			}
		}
		if found {
			continue
		}
		missed = append(missed, tag)
		if _, ok := r.cache.Get(cacheKeyImageInfo(tag)); !ok {
			logrus.WithField("tag", tag).Debug("Tag not found in cache")
			image := fmt.Sprintf("%s/%s:%s", strings.TrimPrefix(r.registry, "https://"), r.repo, tag)
			if err := sem.Acquire(r.ctx, 1); err != nil {
				logrus.WithError(err).WithField("tag", tag).Warn("Failed to acquire semaphore")
				continue
			}
			go func(i string) {
				defer sem.Release(1)
				logrus.WithField("tag", tag).Debug("Fetching metadata")
				info, err := getImageInfo(i)
				if err != nil {
					logrus.WithError(err).WithField("image", image).Warn("Failed to fetch image info")
					return
				}
				r.cache.Set(cacheKeyImageInfo(tag), info, cache.NoExpiration)
			}(image)
		}
	}
	if err := sem.Acquire(r.ctx, int64(r.maxConcurrency)); err != nil {
		logrus.WithError(err).Warn("Failed to acquire semaphore")
	}

	for _, tag := range missed {
		if value, ok := r.cache.Get(cacheKeyImageInfo(tag)); ok {
			logrus.WithField("tag", tag).Debug("Adding a missing tag into the graph")
			info := value.(ImageInfo)
			v, err := semver.Make(info.Version)
			if err != nil {
				logrus.WithError(err).WithField("tag", tag).WithField("version", info.Version).Warn("Failed to parse info.version for tag (ignored)")
				continue
			}
			// TODO: define the func on the struct
			EnsureNode(&graph, nodeWithImageInfo(r.registry, r.repo, tag, v, info))
		} else {
			logrus.WithField("tag", tag).Warn("Tag not found in cache (ignored until the next try)")
		}
	}

	for _, node := range graph.Nodes {
		edges := node.getPrevious(graph)
		// TODO: define the func on the struct
		EnsureEdges(&graph, edges)
	}
	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraped the repository for nodes and edges")
	return graph, nil
}

func EnsureNode(graph *Graph, n Node) {
	for i, node := range graph.Nodes {
		if node.Tag == n.Tag {
			graph.Nodes[i] = n
		}
	}
	graph.Nodes = append(graph.Nodes, n)
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

func nodeWithImageInfo(registry, repo, tag string, version semver.Version, info ImageInfo) Node {
	node := Node{
		Version: version,
		Image:   fmt.Sprintf("%s/%s@%s", strings.TrimPrefix(registry, "https://"), repo, info.Digest),
		Metadata: map[string]string{
			"io.openshift.upgrades.graph.previous.remove_regex": "todo",
			MetadataKeyManifestRef:                              info.Digest,
		},
		Tag:      tag,
		Previous: info.CincinnatiMetadata.Previous,
	}
	if url := info.CincinnatiMetadata.Metadata["url"]; url != "" {
		node.AddMetadata("url", url)
	}
	return node
}

type CincinnatiMetadata struct {
	Kind     string            `json:"kind"`
	Version  string            `json:"version"`
	Previous []string          `json:"previous"`
	Metadata map[string]string `json:"metadata"`
}

type ImageInfo struct {
	CincinnatiMetadata
	Digest string
}

func cacheKeyImageInfo(tag string) string {
	return fmt.Sprintf("imageInfo-%s", tag)
}

func getImageInfo(image string) (ImageInfo, error) {
	var ret ImageInfo
	ref, err := name.ParseReference(image)
	if err != nil {
		return ret, err
	}

	img, err := remote.Image(ref)
	if err != nil {
		return ret, err
	}

	hash, err := img.Digest()
	if err != nil {
		return ret, err
	}
	digest := hash.String()

	layers, err := img.Layers()
	if err != nil {
		return ret, err
	}

	target := "release-manifests/release-metadata"

	for _, layer := range layers {
		rc, err := layer.Uncompressed()
		if err != nil {
			return ret, err
		}

		tr := tar.NewReader(rc)

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return ret, err
			}

			if hdr.Name == target {
				data, err := io.ReadAll(tr)
				if err != nil {
					return ret, err
				}
				var m CincinnatiMetadata
				if err := json.Unmarshal(data, &m); err != nil {
					return ret, err
				}
				ret.Digest = digest
				ret.CincinnatiMetadata = m
				return ret, nil
			}
		}
	}

	return ret, fmt.Errorf("no metadata found for image %s", image)
}
