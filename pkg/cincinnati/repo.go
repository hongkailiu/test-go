package cincinnati

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
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
	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/hongkailiu/test-go/pkg/util"
)

type Repo struct {
	client         Client
	registry       string
	repo           string
	mockDir        string
	dataDir        string
	maxConcurrency int

	releaseMode bool
}

// NewRepo returns a repo.
func NewRepo(client Client, registry, repo, dataDir, mockDir string, maxConcurrency int, releaseMode bool) *Repo {
	return &Repo{
		client:         client,
		registry:       registry,
		repo:           repo,
		dataDir:        dataDir,
		mockDir:        mockDir,
		maxConcurrency: maxConcurrency,
		releaseMode:    releaseMode,
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

func (r *Repo) tags(ctx context.Context) ([]string, error) {
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

	var stop bool
	for !stop {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if !r.releaseMode && count > 2 {
				stop = true
				logrus.WithField("count", count).Info("Break the loop earlier in non-release mode")
				break
			}

			tags, next, err := fetchTags(r.client, url)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch tags: %w", err)
			}

			ret = append(ret, tags...)

			if next != "" {
				nextURL, err := getNextURL(next)
				if err != nil {
					return nil, fmt.Errorf("failed to get next URL: %w", err)
				}

				url = fmt.Sprintf("%s%s", r.registry, nextURL)
			} else {
				stop = true
				break
			}
			count++
		}
	}
	return sets.List[string](sets.New[string](ret...)), nil
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
	logrus.WithField("url", url).Debug("Fetching tags ...")

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
			var e *IsManifestListError
			if ok := errors.As(err, &e); ok {
				logJ.WithError(err).Debug("Failed to fetch image info (will be handled later)")
			} else {
				logJ.WithError(err).Warn("Failed to fetch image info")
			}

			results <- result{err: fmt.Errorf("failed to get image info for tag %s and image %s: %w", job.tag, job.image, err)}

			continue
		}

		logJ.Debug("Fetched image info successfully")

		info.Tag = job.tag
		results <- result{info: info}
	}
}

var invalidTags = sets.New[string](
	// tags whose previous are not smaller
	"4.20.0-ec.7-aarch64", "4.20.0-ec.7-ppc64le", "4.20.0-ec.7-s390x", "4.20.0-ec.7-x86_64", "4.4.0-s390x",
	// multi-tag without the x86_64 variant
	"4.20.0-ec.7-multi",
	// not semver
	"v4.0-20180928191152", "4.11-art-latest-multi",
	// not really multi tags
	"4.7.32-multi", "4.7.32-multi-amd64", "4.9.0-rc.8-multi", "4.9.0-rc.8-multi-amd64", "4.9.0-rc.8-multi-arm64",
)

func (r *Repo) tagsToNodesAndEdges(ctx context.Context, graph Graph) (Graph, error) {

	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraping the repository for nodes and edges ...")

	tags, err := r.tags(ctx)
	if err != nil {
		return Graph{}, fmt.Errorf("failed to fetch tags: %w", err)
	}
	logrus.WithField("tags", len(tags)).Info("Got tags")
	go func(dir string, tags []string) {
		if !r.releaseMode {
			return
		}
		data, err := yaml.Marshal(map[string][]string{"tags": tags})
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal tags")
		}
		file := filepath.Join(dir, "tags.yaml.gz")
		if err := util.WriteBytesMaybeGZIP(file, data); err != nil {
			logrus.WithError(err).WithField("file", file).Error("Failed to write tags to file")
		}
	}(r.dataDir, tags)

	var missing []string
	var multi []ImageInfo
	var invalid int
	for _, tag := range tags {

		if invalid%2000 == 1 {
			logrus.WithField("tag", tag).
				WithField("invalid", invalid).
				WithField("tags", len(tags)).
				Info("Ignored invalid tags")
		}

		multiSuffix := string(ArchTagSuffixMULTI)
		if (strings.HasPrefix(tag, "sha256-") && strings.HasSuffix(tag, ".sig")) ||
			strings.HasSuffix(tag, multiSuffix+"-"+string(ArchTagSuffixAMD64)) ||
			strings.HasSuffix(tag, multiSuffix+"-"+string(ArchTagSuffixARM64)) ||
			strings.HasSuffix(tag, multiSuffix+"-"+string(ArchTagSuffixS390x)) ||
			strings.HasSuffix(tag, multiSuffix+"-"+string(ArchTagSuffixPPC64LE)) ||
			strings.Contains(tag, "nightly") || strings.Contains(tag, "assembly") ||
			invalidTags.Has(tag) {
			logrus.WithField("tag", tag).WithField("invalid", invalid).Debug("Ignored an invalid tag")
			invalid++
			continue
		}
		file := tagToFile(r.dataDir, tag)
		if file == "" {
			logrus.WithField("tag", tag).WithField("invalid", invalid).Debug("Ignored an invalid tag")
			invalid++
			continue
		}

		if graph.Find(tag) > -1 {
			logrus.WithField("tag", tag).Debug("Ignored fetching metadata for an existing tag")
			metrics.tagScraped.WithLabelValues("local").Inc()
			continue
		}

		if fileExists(file) {
			data, err := os.ReadFile(file)
			if err != nil {
				logrus.WithError(err).WithField("tag", tag).WithField("file", file).
					Warn("Failed to fetch image info from file, ignored the tag")
				continue
			}
			info := ImageInfo{}
			if err := yaml.Unmarshal(data, &info); err != nil {
				logrus.WithError(err).WithField("tag", tag).Warn("Failed to unmarshal image info, ignored the tag")
				continue
			}

			if strings.HasSuffix(tag, multiSuffix) && info.Version == "" && info.Digest != "" {
				multi = append(multi, info)
				logrus.WithField("tag", tag).Debug("Ignored a multi tag")
				continue
			}

			version, err := semver.Make(info.Version)
			if err != nil {
				logrus.WithError(err).WithField("tag", info.Tag).WithField("version", info.Version).
					Warn("Failed to parse info.version for tag (ignored and will re-fetch)")

				continue
			}

			graph = graph.EnsureNode(nodeWithImageInfo(r.registry, r.repo, info.Tag, version, info))
			metrics.tagScraped.WithLabelValues("local").Inc()
			continue
		}

		missing = append(missing, tag)
	}

	logrus.WithField("invalid", invalid).
		WithField("tags", len(tags)).
		Info("Ignored invalid tags")

	results := make(chan result, len(missing))
	jobs := make(chan job, len(missing))

	var waitGroup sync.WaitGroup
	for i := 1; i <= r.maxConcurrency; i++ {
		waitGroup.Add(1)

		go worker(i, jobs, results, &waitGroup)
	}

	logrus.WithField("total", len(missing)).Info("Fetching image metadata for missing tags ...")

	for i, tag := range missing {
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

	var received int

	for i := 0; i < len(missing); i++ {
		select {
		case result := <-results:
			logrus.WithField("received", received).WithField("total", len(missing)).Debug("Received result")
			received++

			if result.err != nil {
				var e *IsManifestListError
				if ok := errors.As(result.err, &e); ok {
					i := strings.LastIndex(e.Image, ":")
					if i == -1 {
						logrus.WithField("image", e.Image).Warn("Ignored an invalid multi image")
						continue
					}
					tag := e.Image[i+1:]
					// The other fields will be filled with the image info from its amd64 shard
					info := ImageInfo{
						Digest: e.Digest,
						Tag:    tag,
					}
					multi = append(multi, info)
					go func(dir string, info ImageInfo) {
						saveToFile(dir, info)
					}(r.dataDir, info)
					logrus.WithField("tag", tag).Debug("Ignored a multi tag in a received result (will be handled later)")
					continue
				}

				logrus.WithError(result.err).Warn("Failed to get image info and the tag is ignored")

				continue
			}

			info := result.info
			logrus.WithField("tag", info.Tag).Debug("Adding a missing tag into the graph")

			go func(dir string, info ImageInfo) {
				saveToFile(dir, info)
			}(r.dataDir, info)

			version, err := semver.Make(info.Version)
			if err != nil {
				logrus.WithError(err).WithField("tag", info.Tag).WithField("version", info.Version).
					Warn("Failed to parse info.version for tag (ignored)")

				continue
			}

			graph = graph.EnsureNode(nodeWithImageInfo(r.registry, r.repo, info.Tag, version, info))
			nodes++
			metrics.tagScraped.WithLabelValues("upstream").Inc()
			if nodes%100 == 1 {
				logrus.WithField("tag", info.Tag).
					WithField("missing", len(missing)).
					WithField("received", received).
					WithField("nodes", nodes).
					WithField("tags", len(tags)).
					Info("Added nodes to the graph")
			}
		case <-ctx.Done():
			logrus.Info("Context done")
			return graph, ctx.Err()
		}
	}

	logrus.WithField("missing", len(missing)).
		WithField("received", received).
		WithField("nodes", nodes).
		WithField("tags", len(tags)).
		Info("Finished scraping the repository graph")

	logrus.WithField("nodes", len(graph.Nodes)).WithField("multi", len(multi)).Info("Adding multi nodes the repository graph ...")
	for _, info := range multi {
		tag := info.Tag
		trimmed := strings.TrimSuffix(tag, string(ArchTagSuffixMULTI))
		if trimmed == tag {
			logrus.WithField("tag", tag).Warn("Ignored an invalid multi image")
			continue
		}
		// The other fields are filled with the image info from its amd64 shard
		// `oc adm release info` does a similar thing. By default, it takes the arch on machine where `oc` runs.
		i := graph.Find(trimmed + string(ArchTagSuffixAMD64))
		if i == -1 {
			logrus.WithField("tag", tag).Warn("Failed to find the amd64 node, ignored an invalid multi image")
			continue
		}

		node := Node{
			Version: graph.Nodes[i].Version,
			Image:   tagToImage(r.registry, r.repo, info.Digest),
			Tag:     tag,
		}
		node.Metadata = map[string]string{}
		SetMetadata(&node, MetadataKeyArchitecture, "multi")
		if graph.Nodes[i].Metadata != nil {
			for k, v := range graph.Nodes[i].Metadata {
				node.Metadata[k] = v
			}
		}
		if graph.Nodes[i].Previous != nil {
			node.Previous = append(node.Previous, graph.Nodes[i].Previous...)
		}

		graph = graph.EnsureNode(node)
	}

	logrus.WithField("nodes", len(graph.Nodes)).WithField("multi", len(multi)).Info("Added multi nodes the repository graph")

	for _, node := range graph.Nodes {
		edges := node.getPrevious(graph)
		graph = graph.EnsureEdges(edges)
	}

	logrus.WithField("nodes", len(graph.Nodes)).WithField("edges", len(graph.Edges)).WithField("conditionalEdges", len(graph.ConditionalEdges)).
		Info("Scraped the repository for nodes and edges")

	return graph, nil
}

func saveToFile(dir string, info ImageInfo) {
	logger := logrus.WithField("tag", info.Tag)

	file := tagToFile(dir, info.Tag)
	if file == "" {
		logger.Warn("Failed to determine the file, skipped")
		return
	}
	if fileExists(file) {
		logger.WithField("file", file).Warn("File already exists, skipped")
		return
	}
	logger = logger.WithField("file", file)
	logger.Debug("Saving to disk ...")
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		logger.WithError(err).Warn("Failed to create directory")
		return
	}
	data, err := yaml.Marshal(info)
	if err != nil {
		logger.WithError(err).Warn("Failed to marshal info")
		return
	}
	err = os.WriteFile(file, data, 0600)
	if err != nil {
		logger.WithError(err).WithField("file", file).Warn("Failed to write file")
		return
	}
	logger.Debug("Saved to disk ...")
}

func tagToFile(dir, tag string) string {
	vTag := tag
	i := strings.Index(tag, "-")
	if i != -1 {
		vTag = tag[:i]
	}
	version, err := semver.Parse(vTag)
	if err != nil {
		logrus.WithField("tag", tag).WithField("vTag", vTag).Warn("Failed to parse version from tag")
		return ""
	}
	return filepath.Join(dir, fmt.Sprintf("%d.%d", version.Major, version.Minor), tag+".yaml")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	// The file might exist but be a directory; return true only for regular files
	return err == nil && info.Mode().IsRegular()
}

func tagToImage(registry, repo, digest string) string {
	return fmt.Sprintf("%s/%s@%s", strings.TrimPrefix(registry, "https://"), repo, digest)
}

func nodeWithImageInfo(registry, repo, tag string, version semver.Version, info ImageInfo) Node {
	node := Node{
		Version:  version,
		Image:    tagToImage(registry, repo, info.Digest),
		Tag:      tag,
		Previous: info.Previous,
	}
	if url, ok := info.Metadata["url"]; ok {
		SetMetadata(&node, "url", url)
	}

	if arch, ok := info.Metadata[MetadataKeyArchitecture]; ok {
		SetMetadata(&node, MetadataKeyArchitecture, arch)
	}

	return node
}

type CincinnatiMetadata struct {
	Kind     string            `json:"kind"`
	Version  string            `json:"version"`
	Previous []string          `json:"previous"`
	Metadata map[string]string `json:"metadata"`
}

type IsManifestListError struct {
	Image  string
	Digest string
}

func (e *IsManifestListError) Error() string {
	return fmt.Sprintf("image %s is a manifest list", e.Image)
}

type ImageInfo struct {
	CincinnatiMetadata

	Digest string

	Tag string
}

func getImageInfo(image string) (ImageInfo, error) {
	var ret ImageInfo

	ref, err := name.ParseReference(image)
	if err != nil {
		return ret, fmt.Errorf("failed to parse reference of image %s: %w", image, err)
	}

	desc, err := remote.Get(ref)
	if err != nil {
		return ret, fmt.Errorf("failed to get descriptor for %s: %w", image, err)
	}

	digest := desc.Digest.String()
	if digest == "" {
		return ret, fmt.Errorf("failed to get digest for %s", image)
	}

	// Check if it's a manifest list (index)
	if desc.MediaType.IsIndex() {
		// A sanity check. A contract with ART that multi-arch image always with a tag with suffix "multi".
		// If the contract breaks, Cincinnati ignores those payload images in the update graph.
		multiSuffix := string(ArchTagSuffixMULTI)
		if !strings.HasSuffix(image, multiSuffix) {
			return ret, fmt.Errorf("multi-arch image %s does not end with %s", image, ArchTagSuffixMULTI)
		}
		// The image info will be recovered with the image info from its amd64 shard
		return ret, &IsManifestListError{
			Image:  image,
			Digest: digest,
		}
	}

	img, err := remote.Image(ref)
	if err != nil {
		return ret, fmt.Errorf("failed to fetch image %s: %w", image, err)
	}

	layers, err := img.Layers()
	if err != nil {
		return ret, fmt.Errorf("error getting layers for image %q: %w", image, err)
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
				return ret, fmt.Errorf("error getting the next tar archive: %w", err)
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
