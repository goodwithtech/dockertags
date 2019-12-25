package utils

import (
	"strings"

	"github.com/goodwithtech/dockertags/internal/types"
)

// MatchConditionTags retunrn matched option
func MatchConditionTags(opt *types.FilterOption, tagNames []string) (contained bool) {
	if len(opt.Contain) == 0 {
		return true
	}
	for _, tagName := range tagNames {
		var currentTagContained int
		for _, target := range opt.Contain {
			if strings.Contains(tagName, target) {
				currentTagContained++
			}
		}
		if len(opt.Contain) == currentTagContained {
			return true
		}
	}
	return false
}
