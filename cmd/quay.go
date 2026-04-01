package cmd

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type CincinnatiMetadata struct {
	Kind     string            `json:"kind"`
	Version  string            `json:"version"`
	Previous []string          `json:"previous"`
	Metadata map[string]string `json:"metadata"`
}

func getCincinnatiMetadata(image string) (CincinnatiMetadata, error) {
	var ret CincinnatiMetadata

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

				if err := json.Unmarshal(data, &ret); err != nil {
					return ret, fmt.Errorf("failed to unmarshal image info: %w", err)
				}

				return ret, nil
			}
		}
	}

	return ret, fmt.Errorf("no metadata found for image %s", image)
}
