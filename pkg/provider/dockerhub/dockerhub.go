package dockerhub

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/goodwithtech/dockertags/internal/log"
	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"

	dockertypes "github.com/docker/docker/api/types"
)

const registryURL = "https://registry.hub.docker.com/"

type DockerHub struct {
	filterOpt *types.FilterOption
}

type tagsResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []ImageSummary `json:"results"`
}

type ImageSummary struct {
	Name        string `json:"name"`
	FullSize    int    `json:"full_size"`
	LastUpdated string `json:"last_updated"`
}

func (p *DockerHub) Run(ctx context.Context, domain, repository string, reqOpt *types.RequestOption, filterOpt *types.FilterOption) (types.ImageTags, error) {
	p.filterOpt = filterOpt
	auth := dockertypes.AuthConfig{
		ServerAddress: "registry.hub.docker.com",
		Username:      reqOpt.UserName,
		Password:      reqOpt.Password,
	}
	// fetch page 1 for check max item count.
	tagResp, err := getTagResponse(ctx, auth, reqOpt.Timeout, repository, 1)
	if err != nil {
		return nil, err
	}
	imageTags := p.convertResultToTag(tagResp.Results)
	if reqOpt.MaxCount > 0 && len(imageTags) > reqOpt.MaxCount {
		return imageTags, nil
	}

	lastPage := calcMaxRequestPage(tagResp.Count, reqOpt.MaxCount, filterOpt)
	// create ch (page - 1), already fetched first page,
	tagsPerPage := make(chan types.ImageTags, lastPage-1)
	eg := errgroup.Group{}
	for page := 2; page <= lastPage; page++ {
		page := page
		eg.Go(func() error {
			tagResp, err := getTagResponse(ctx, auth, reqOpt.Timeout, repository, page)
			if err != nil {
				return err
			}
			tagsPerPage <- p.convertResultToTag(tagResp.Results)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for page := 2; page <= lastPage; page++ {
		select {
		case tags := <-tagsPerPage:
			imageTags = append(imageTags, tags...)
		}
	}
	return imageTags, nil
}

func (p *DockerHub) convertResultToTag(summaries []ImageSummary) types.ImageTags {
	tags := []types.ImageTag{}
	for _, detail := range summaries {
		if detail.Name == "" {
			continue
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, detail.LastUpdated)
		tagNames := []string{detail.Name}
		if !utils.MatchConditionTags(p.filterOpt, tagNames) {
			continue
		}
		tags = append(tags, types.ImageTag{
			Tags:      tagNames,
			Byte:      &detail.FullSize,
			CreatedAt: &createdAt,
		})
	}
	return tags
}

// getTagResponse returns the tags for a specific repository.
// curl 'https://registry.hub.docker.com/v2/repositories/library/debian/tags/'
func getTagResponse(ctx context.Context, auth dockertypes.AuthConfig, timeout time.Duration, repository string, page int) (tagsResponse, error) {
	url := fmt.Sprintf("%s/v2/repositories/%s/tags/?page=%d", registryURL, repository, page)
	log.Logger.Debugf("url=%s,repository=%s", url, repository)
	var response tagsResponse
	if _, err := getJSON(ctx, url, auth, timeout, &response); err != nil {
		return response, err
	}

	return response, nil
}

func calcMaxRequestPage(totalCnt, needCnt int, option *types.FilterOption) int {
	maxPage := totalCnt/types.ITEM_PER_PAGE + 1
	if needCnt == 0 || option.Contain != "" {
		return maxPage
	}
	needPage := needCnt/types.ITEM_PER_PAGE + 1
	if needPage >= maxPage {
		return maxPage
	}
	return needPage
}
