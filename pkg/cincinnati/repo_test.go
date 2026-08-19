package cincinnati

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/go-cmp/cmp"

	"github.com/hongkailiu/test-go/pkg/testhelper"
)

type fakeClient struct {
	cancel context.CancelFunc
	doneAt int
	count  int
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	data := TagsListData{

		Tags: []string{strconv.Itoa(c.count * 2), strconv.Itoa(c.count*2 + 1)},
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling tags: %w", err)
	}

	ret := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(raw)),
		Header:     make(http.Header),
	}

	if c.count < 6 {
		ret.Header.Set("Link", `</v2/openshift-release-dev/ocp-release/tags/list?n=100&last=4.10.36-aarch64>; rel="next"`)
	}

	if c.count == c.doneAt {
		c.cancel()
	}
	c.count++
	return ret, nil
}

func TestRepo_tags(t *testing.T) {
	tests := []struct {
		name        string
		doneAt      int
		releaseMode bool
		expTags     []string
		expErr      error
	}{
		{
			name:   "context canceled",
			doneAt: 2,
			expErr: context.Canceled,
		},
		{
			name:        "context canceled in release mode",
			doneAt:      3,
			expErr:      context.Canceled,
			releaseMode: true,
		},
		{
			name:    "non-release mode",
			doneAt:  3,
			expTags: []string{"0", "1", "2", "3", "4", "5"},
		},
		{
			name:        "complete",
			doneAt:      -1,
			expTags:     []string{"0", "1", "10", "11", "12", "13", "2", "3", "4", "5", "6", "7", "8", "9"},
			releaseMode: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			repo := Repo{
				client:      &fakeClient{doneAt: test.doneAt, cancel: cancel},
				releaseMode: test.releaseMode,
			}
			actual, actualErr := repo.tags(ctx)
			if diff := cmp.Diff(test.expErr, actualErr, testhelper.ErrorTransformer); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.expTags, actual); diff != "" {
				t.Errorf("tags mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetImageInfo(t *testing.T) {
	tests := []struct {
		name           string
		image          string
		want           ImageInfo
		expErrContains string
	}{
		{
			name:  "okd scos release",
			image: "quay.io/okd/scos-release:4.22.0-okd-scos.7",
			want: ImageInfo{
				CincinnatArchitecture: Multi,
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:    "cincinnati-metadata-v0",
					Version: "4.22.0-okd-scos.7",
					Previous: []string{
						"4.21.0-okd-scos.11",
						"4.22.0-okd-scos.0",
						"4.22.0-okd-scos.1",
						"4.22.0-okd-scos.2",
						"4.22.0-okd-scos.3",
						"4.22.0-okd-scos.4",
						"4.22.0-okd-scos.5",
						"4.22.0-okd-scos.6",
					},
				},
				Digest:       "sha256:790386476b24e1a1043aa865c98d2bfeebdd467aa54859f0f53ecd6c487a8036",
				Tag:          "4.22.0-okd-scos.7",
				OS:           "linux",
				Architecture: "amd64",
				Manifests: []Manifest{
					{
						OS:           "linux",
						Architecture: "amd64",
						Digest:       "sha256:790386476b24e1a1043aa865c98d2bfeebdd467aa54859f0f53ecd6c487a8036",
					},
					{
						OS:           "linux",
						Architecture: "arm64",
						Digest:       "sha256:27a7ffdc2a7b6a03cb2969fc6cf8609c560d85eb3c3af1ee4bfce90c560af16f",
					},
				},
				IsIndex: true,
			},
		},
		{
			name:  "ocp release multi",
			image: "quay.io/openshift-release-dev/ocp-release:4.15.5-multi",
			want: ImageInfo{
				CincinnatArchitecture: Multi,
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:    "cincinnati-metadata-v0",
					Version: "4.15.5",
					Previous: []string{
						"4.14.14",
						"4.14.15",
						"4.14.16",
						"4.14.17",
						"4.14.18",
						"4.15.0",
						"4.15.2",
						"4.15.3",
					},
					Metadata: map[string]string{
						MetadataKeyArchitecture: Multi,
						"url":                   "https://access.redhat.com/errata/RHSA-2024:1449",
					},
				},
				Digest:       "sha256:c07f82ae6a59893a4a3c5f2780fbe931382a737f23fe580fb3d2a07b7fd011f5",
				Tag:          "4.15.5-multi",
				OS:           "linux",
				Architecture: "amd64",
				Manifests: []Manifest{
					{
						OS:           "linux",
						Architecture: "amd64",
						Digest:       "sha256:c07f82ae6a59893a4a3c5f2780fbe931382a737f23fe580fb3d2a07b7fd011f5",
					},
					{
						OS:           "linux",
						Architecture: "ppc64le",
						Digest:       "sha256:4a7c0f888b0a89157d7e9f48df55a361a9bfc009f63be560de79c24aa8acf8d0",
					},
					{
						OS:           "linux",
						Architecture: "s390x",
						Digest:       "sha256:f5365e8eda348e162bdadac46bbefdf419f5d50a07db8a5b1cfdf96ef50bf159",
					},
					{
						OS:           "linux",
						Architecture: "arm64",
						Digest:       "sha256:15a384b58264c3c8f938cb258cfffb9abc9561775db0d0585caeeaf41e79bd8e",
					},
				},
				IsIndex: true,
			},
		},
		{
			name:  "okd scos release single arch index",
			image: "quay.io/okd/scos-release:4.16.0-okd-scos.0",
			want: ImageInfo{
				CincinnatArchitecture: "amd64",
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:     "cincinnati-metadata-v0",
					Version:  "4.16.0-okd-scos.0",
					Previous: []string{},
				},
				Digest:       "sha256:830fef3e885e0eb63b2a4fc5352cb1729506d5ec3ba21c71fd26247c3530312d",
				Tag:          "4.16.0-okd-scos.0",
				OS:           "linux",
				Architecture: "amd64",
				Manifests: []Manifest{
					{
						OS:           "linux",
						Architecture: "amd64",
						Digest:       "sha256:830fef3e885e0eb63b2a4fc5352cb1729506d5ec3ba21c71fd26247c3530312d",
					},
				},
				IsIndex: true,
			},
		},
		{
			name:  "ocp release by digest (4.15.5-multi)",
			image: "quay.io/openshift-release-dev/ocp-release@sha256:c07f82ae6a59893a4a3c5f2780fbe931382a737f23fe580fb3d2a07b7fd011f5",
			want: ImageInfo{
				CincinnatArchitecture: Multi,
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:    "cincinnati-metadata-v0",
					Version: "4.15.5",
					Previous: []string{
						"4.14.14",
						"4.14.15",
						"4.14.16",
						"4.14.17",
						"4.14.18",
						"4.15.0",
						"4.15.2",
						"4.15.3",
					},
					Metadata: map[string]string{
						MetadataKeyArchitecture: Multi,
						"url":                   "https://access.redhat.com/errata/RHSA-2024:1449",
					},
				},
				Digest:       "sha256:c07f82ae6a59893a4a3c5f2780fbe931382a737f23fe580fb3d2a07b7fd011f5",
				OS:           "linux",
				Architecture: "amd64",
				IsIndex:      false,
			},
		},
		{
			name:  "ocp release multi s390x tag",
			image: "quay.io/openshift-release-dev/ocp-release:4.22.9-multi-s390x",
			want: ImageInfo{
				CincinnatArchitecture: Multi,
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:    "cincinnati-metadata-v0",
					Version: "4.22.9",
					Previous: []string{
						"4.21.10",
						"4.21.11",
						"4.21.12",
						"4.21.13",
						"4.21.14",
						"4.21.15",
						"4.21.16",
						"4.21.17",
						"4.21.18",
						"4.21.19",
						"4.21.20",
						"4.21.21",
						"4.21.22",
						"4.21.23",
						"4.21.24",
						"4.21.25",
						"4.21.26",
						"4.21.27",
						"4.21.28",
						"4.21.6",
						"4.21.7",
						"4.21.8",
						"4.21.9",
						"4.22.0",
						"4.22.1",
						"4.22.2",
						"4.22.3",
						"4.22.4",
						"4.22.5",
						"4.22.6",
						"4.22.7",
						"4.22.8",
					},
					Metadata: map[string]string{
						MetadataKeyArchitecture: Multi,
						"url":                   "https://access.redhat.com/errata/RHSA-2026:51038",
					},
				},
				Digest:       "sha256:9a3565331bf526bd2f4bbac974fa6e19e860f26b0b2f061378efb7191c7acbbf",
				Tag:          "4.22.9-multi-s390x",
				OS:           "linux",
				Architecture: "s390x",
				IsIndex:      false,
			},
		},
		{
			name:  "okd scos release ec.3",
			image: "quay.io/okd/scos-release:5.0.0-okd-scos.ec.3",
			want: ImageInfo{
				CincinnatArchitecture: "amd64",
				CincinnatiMetadata: CincinnatiMetadata{
					Kind:    "cincinnati-metadata-v0",
					Version: "5.0.0-okd-scos.ec.3",
					Previous: []string{
						"4.22.0-okd-scos.ec.16",
						"5.0.0-okd-scos.ec.0",
						"5.0.0-okd-scos.ec.1",
						"5.0.0-okd-scos.ec.2",
					},
				},
				Digest:       "sha256:a4f0bfa8834dd89d9eed50083c43298b2b1066a587734090f67458fbd067d31e",
				Tag:          "5.0.0-okd-scos.ec.3",
				OS:           "linux",
				Architecture: "amd64",
				Manifests: []Manifest{
					{
						OS:           "linux",
						Architecture: "amd64",
						Digest:       "sha256:a4f0bfa8834dd89d9eed50083c43298b2b1066a587734090f67458fbd067d31e",
					},
				},
				IsIndex: true,
			},
		},
		{
			name:           "invalid reference",
			image:          "!!!not-a-valid-reference!!!",
			expErrContains: "failed to parse reference",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetImageInfo(test.image)
			if test.expErrContains != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if !strings.Contains(err.Error(), test.expErrContains) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ImageInfo mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRepo_tagsToNodesAndEdges(t *testing.T) {
	tests := []struct {
		name   string
		exp    Graph
		expErr error
	}{
		{
			name: "mock",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := Repo{
				client:         http.DefaultClient,
				dataDir:        "../../mock/data",
				registry:       "https://quay.io",
				repo:           "openshift-release-dev/ocp-release",
				maxConcurrency: 1,
			}
			var graph Graph
			actual, actualErr := repo.tagsToNodesAndEdges(context.Background(), graph)
			if diff := cmp.Diff(test.expErr, actualErr, testhelper.ErrorTransformer); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}

			if actualErr != nil {
				return
			}

			data, err := json.MarshalIndent(actual, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			cupaloy.New(
				cupaloy.SnapshotSubdirectory("testdata/.snapshots"),
			).SnapshotT(t, data)
		})
	}
}
