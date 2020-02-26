package dockerhub

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/goodwithtech/dockertags/internal/types"
)

func TestScanImage(t *testing.T) {
	testcases := map[string]struct {
		filePath  string
		filterOpt types.FilterOption
		expected  types.ImageTags
	}{
		"debian page1": {
			filePath:  "./testdata/page1.json",
			filterOpt: types.FilterOption{},
			expected: types.ImageTags{
				types.ImageTag{
					Tags: []string{"unstable-slim", "unstable-20191118-slim"},
				},
				types.ImageTag{
					Tags: []string{"unstable", "unstable-20191118"},
				},
				types.ImageTag{
					Tags: []string{"testing-slim", "testing-20191118-slim"},
				},
				types.ImageTag{
					Tags: []string{"testing-backports"},
				},
				types.ImageTag{
					Tags: []string{"testing", "testing-20191118"},
				},
				types.ImageTag{
					Tags: []string{"stretch-slim"},
				},
			},
		},
		"debian filter slim": {
			filePath:  "./testdata/page1.json",
			filterOpt: types.FilterOption{Contain: []string{"slim"}},
			expected: types.ImageTags{
				types.ImageTag{
					Tags: []string{"unstable-slim", "unstable-20191118-slim"},
				},
				types.ImageTag{
					Tags: []string{"testing-slim", "testing-20191118-slim"},
				},
				types.ImageTag{
					Tags: []string{"stretch-slim"},
				},
			},
		},
	}

	for tc, v := range testcases {
		dockerHub := DockerHub{filterOpt: &v.filterOpt}
		var data tagsResponse
		file, err := ioutil.ReadFile(v.filePath)
		if err != nil {
			t.Errorf("readfile error: %w", err)
			continue
		}
		json.Unmarshal(file, &data)
		actual := dockerHub.convertResultToTag(data.Results)
		opts := []cmp.Option{
			cmp.Transformer("Sort", func(in []string) []string {
				out := append([]string{}, in...) // Copy input to avoid mutating it
				sort.Strings(out)
				return out
			}),
			cmpopts.IgnoreFields(types.ImageTag{}, "Byte", "UploadedAt"),
		}
		sort.Sort(actual)
		if diff := cmp.Diff(v.expected, actual, opts...); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}

func TestSummarizeByHash(t *testing.T) {
	testcases := map[string]struct {
		tags     []tagSummary
		expected map[string]types.ImageTag
	}{
		"OK": {
			tags: []tagSummary{
				{
					Name:        "a",
					LastUpdated: "2019-12-02T00:00:00.00000Z",
					Images: []image{
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "400", Architecture: "995test"},
					},
				},
				{
					Name:        "b",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400b", Architecture: "995test"},
						{Digest: "001b", Architecture: "999test"},
						{Digest: "100b", Architecture: "998test"},
						{Digest: "200b", Architecture: "997test"},
						{Digest: "300b", Architecture: "996test"},
					},
				},
				{
					Name:        "c",
					LastUpdated: "2019-12-03T00:00:00.00000Z",
					Images: []image{
						{Digest: "400", Architecture: "995test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
					},
				},
			},
			expected: map[string]types.ImageTag{
				"400": {
					Tags:       []string{"a", "c"},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"400b": {
					Tags:       []string{"b"},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"LoadUpdatedAt": {
			tags: []tagSummary{
				{
					Name:        "a",
					LastUpdated: "2019-12-02T00:00:00.00000Z",
					Images: []image{
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "400", Architecture: "995test"},
					},
				},
				{
					Name:        "b",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400b", Architecture: "995test"},
						{Digest: "001b", Architecture: "999test"},
						{Digest: "100b", Architecture: "998test"},
						{Digest: "200b", Architecture: "997test"},
						{Digest: "300b", Architecture: "996test"},
					},
				},
				{
					Name:        "c",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400", Architecture: "995test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
					},
				},
			},
			expected: map[string]types.ImageTag{
				"400": {
					Tags:       []string{"a", "c"},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"400b": {
					Tags:       []string{"b"},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for tc, v := range testcases {
		actual := summarizeByHash(v.tags)
		if diff := cmp.Diff(v.expected, actual); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}
