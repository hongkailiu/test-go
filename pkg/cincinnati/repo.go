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
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/sets"
)

type Repo struct {
	client         Client
	registry       string
	repo           string
	mockDir        string
	maxConcurrency int
}

func NewRepo(client Client, registry, repo, mockDir string, maxConcurrency int, cache Cache) *Repo {
	return &Repo{
		client:         client,
		registry:       registry,
		repo:           repo,
		mockDir:        mockDir,
		maxConcurrency: maxConcurrency,
	}
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Cache interface {
	Get(k string) (any, bool)
	Set(k string, x any, d time.Duration)
}

type TagsListData struct {
	Tags []string `json:"tags"`
}

func (r *Repo) tags() ([]string, error) {
	if r.mockDir != "" {
		raw, err := os.ReadFile(filepath.Join(r.mockDir, "repo.tags.list.json"))
		if err != nil {
			return nil, err
		}

		data := TagsListData{}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		return data.Tags, nil
	}

	var ret []string

	url := fmt.Sprintf("%s/v2/%s/tags/list", r.registry, r.repo)

	var count int

	notReleaseMode := !releaseMode()

	for {
		count++
		if notReleaseMode && count > 2 {
			break
		}

		tags, next, err := fetchTags(r.client, url)
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

func fetchTags(client Client, url string) ([]string, string, error) {
	// Reference https://oneuptime.com/blog/post/2026-02-08-how-to-list-all-tags-of-a-docker-image-on-docker-hub/view
	// https://quay.io/v2/openshift-release-dev/ocp-release/tags/list
	logrus.WithField("url", url).Info("Fetching tags ...")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request with url %s: %w", url, err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get response to fetch tags: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code %d for url %s", res.StatusCode, req.URL)
	}

	link := res.Header.Get("Link")
	if link != "" {
		logrus.WithField("link", link).Debugf("Found Link")
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading body: %w", err)
	}

	data := TagsListData{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return data.Tags, link, nil
}

type result struct {
	info ImageInfo
	err  error
}
type job struct {
	i     int
	tag   string
	image string
}

func worker(id int, jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := logrus.WithField("worker.id", id)
	for job := range jobs {
		logJ := logger.WithField("image", job.image).WithField("tag", job.tag).WithField("job.index", job.i)
		logJ.Debug("Worker processing job ...")

		info, err := getImageInfo(job.image)
		if err != nil {
			logJ.WithError(err).Warn("Failed to fetch image info")

			results <- result{err: fmt.Errorf("failed to get image info for tag %s and image %s: %w", job.tag, job.image, err)}

			return
		}

		logJ.Debug("Fetched image info successfully")

		info.Tag = job.tag
		results <- result{info: info}
	}
}

func (r *Repo) tagsToNodesAndEdges(_ context.Context, graph Graph) (Graph, error) {
	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraping the repository for nodes and edges ...")

	tags, err := r.tags()
	// TODO: ignore the tags that are not in channels
	if err != nil {
		return Graph{}, fmt.Errorf("failed to fetch tags: %w", err)
	}

	missing := sets.New[string]()
	for _, tag := range tags {
		if graph.Find(tag) > -1 {
			logrus.WithField("tag", tag).Debug("Ignored fetching metadata for an existing tag")
			continue
		}

		missing.Insert(tag)
	}

	results := make(chan result, len(missing))
	jobs := make(chan job, len(missing))

	var waitGroup sync.WaitGroup
	for i := 1; i <= r.maxConcurrency; i++ {
		waitGroup.Add(1)

		go worker(i, jobs, results, &waitGroup)
	}

	logrus.WithField("total", len(missing)).Debug("Fetching image metadata ...")

	for i, tag := range missing.UnsortedList() {
		image := fmt.Sprintf("%s/%s:%s", strings.TrimPrefix(r.registry, "https://"), r.repo, tag)
		logrus.WithField("image", image).WithField("tag", tag).WithField("index", i).Debug("Sending a job")

		jobs <- job{i: i, tag: tag, image: image}
	}

	close(jobs)

	go func() {
		logrus.Debug("Waiting for fetching image info ...")
		waitGroup.Wait()
		close(results)
		logrus.Debug("Done with waiting for fetching image info ...")
	}()

	var nodes int

	for result := range results {
		if result.err != nil {
			logrus.WithError(result.err).Warn("Failed to get image info and the tag is ignored")

			continue
		}

		info := result.info
		logrus.WithField("tag", info.Tag).Debug("Adding a missing tag into the graph")

		version, err := semver.Make(info.Version)
		if err != nil {
			logrus.WithError(err).WithField("tag", info.Tag).WithField("version", info.Version).
				Warn("Failed to parse info.version for tag (ignored)")

			continue
		}

		nodes++
		graph = graph.EnsureNode(nodeWithImageInfo(r.registry, r.repo, info.Tag, version, info))
	}

	logrus.WithField("messing", len(missing)).WithField("nodes", nodes).WithField("tags", len(tags)).
		Debug("Finished scraping the repository graph ...")

	for _, node := range graph.Nodes {
		edges := node.getPrevious(graph)
		graph = graph.EnsureEdges(edges)
	}

	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraped the repository for nodes and edges")

	return graph, nil
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
		Previous: info.Previous,
	}
	if url, ok := info.Metadata["url"]; ok {
		node.AddMetadata("url", url)
	}

	if arch, ok := info.Metadata[MetadataKeyArchitecture]; ok {
		node.AddMetadata(MetadataKeyArchitecture, arch)
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

	Tag string
}

func cacheKeyImageInfo(tag string) string {
	return "imageInfo-" + tag
}

func getImageInfo(image string) (ImageInfo, error) {
	var ret ImageInfo

	ref, err := name.ParseReference(image)
	if err != nil {
		return ret, fmt.Errorf("failed to parse reference of image %s: %w", image, err)
	}

	img, err := remote.Image(ref)
	if err != nil {
		return ret, fmt.Errorf("failed to fetch image %s: %w", image, err)
	}

	hash, err := img.Digest()
	if err != nil {
		return ret, fmt.Errorf("failed to get digest for image %s: %w", image, err)
	}

	digest := hash.String()

	layers, err := img.Layers()
	if err != nil {
		return ret, fmt.Errorf("error getting layers for image %q: %v", image, err)
	}

	target := "release-manifests/release-metadata"

	for _, layer := range layers {
		rc, err := layer.Uncompressed()
		if err != nil {
			return ret, fmt.Errorf("failed to uncompress layer: %w", err)
		}

		tr := tar.NewReader(rc)

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				return ret, fmt.Errorf("error getting the next tar archive: %v", err)
			}

			if hdr.Name == target {
				data, err := io.ReadAll(tr)
				if err != nil {
					return ret, fmt.Errorf("error reading tar archive: %w", err)
				}

				var m CincinnatiMetadata
				if err := json.Unmarshal(data, &m); err != nil {
					return ret, fmt.Errorf("failed to unmarshal image info: %w", err)
				}

				ret.Digest = digest
				ret.CincinnatiMetadata = m

				return ret, nil
			}
		}
	}

	return ret, fmt.Errorf("no metadata found for image %s", image)
}
