package dockerhub

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"testing"

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
			cmpopts.IgnoreFields(types.ImageTag{}, "Byte", "CreatedAt"),
		}
		sort.Sort(actual)
		if diff := cmp.Diff(v.expected, actual, opts...); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}
