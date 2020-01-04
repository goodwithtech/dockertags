package utils

import (
	"testing"

	"github.com/goodwithtech/dockertags/internal/types"
)

func TestMatchConditionTags(t *testing.T) {
	testcases := map[string]struct {
		tags     []string
		filter   types.FilterOption
		expected bool
	}{
		"NotContain": {
			tags: []string{"test1"},
			filter: types.FilterOption{
				Contain: []string{
					"testing",
				},
			},
			expected: false,
		},
		"SingleContainSingleTag": {
			tags: []string{"test1"},
			filter: types.FilterOption{
				Contain: []string{
					"test",
				},
			},
			expected: true,
		},
		"NotContainsBoth": {
			tags: []string{"test1"},
			filter: types.FilterOption{
				Contain: []string{
					"sample", "test",
				},
			},
			expected: false,
		},
		"ContainBoth": {
			tags: []string{"test1"},
			filter: types.FilterOption{
				Contain: []string{
					"1", "test",
				},
			},
			expected: true,
		},
		"Contain1Tag": {
			tags: []string{"test1", "test2"},
			filter: types.FilterOption{
				Contain: []string{
					"1", "test",
				},
			},
			expected: true,
		},
	}

	for tc, v := range testcases {
		actual := MatchConditionTags(&v.filter, v.tags)
		if actual != v.expected {
			t.Errorf("%s: expected %v, got %v", tc, v.expected, actual)
		}
	}
}
