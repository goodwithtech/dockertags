package utils

import (
	"fmt"
	"strings"

	"github.com/goodwithtech/dockertags/internal/types"
)

// MatchConditionTags retunrn matched option
func MatchConditionTags(opt *types.FilterOption, tagNames []string) (contained bool) {
	if opt.Contain == "" {
		return true
	}
	for _, tagName := range tagNames {
		if strings.Contains(tagName, opt.Contain) {
			return true
		}
	}
	return false
}

func GetWithRepositoryName(repo string, tagNames []string) string {
	return fmt.Sprintf("%s:%s", repo, tagNames[0])
}
