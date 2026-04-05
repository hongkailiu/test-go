package cincinnati

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
				mockDir:        "../../mock",
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
