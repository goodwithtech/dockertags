package dockerhub

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/goodwithtech/dockertags/internal/log"
	"github.com/goodwithtech/dockertags/internal/types"

	dockertypes "github.com/docker/docker/api/types"

	"github.com/goodwithtech/dockertags/pkg/registry"
)

const TAG_PER_PAGE = 10
const registryURL = "https://registry.hub.docker.com/"

type DockerHub struct {
	registry *registry.Registry
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

// curl 'https://registry.hub.docker.com/v2/repositories/library/debian/tags/'
func (p *DockerHub) Run(ctx context.Context, domain, repository string, option types.RequestOption) (types.ImageTags, error) {
	auth := dockertypes.AuthConfig{
		ServerAddress: "registry.hub.docker.com",
		Username:      option.UserName,
		Password:      option.Password,
	}
	// 1ページ目は普通に取得
	tagResp, err := getTagResponse(ctx, auth, option.Timeout, repository, 1)
	if err != nil {
		return nil, err
	}
	imageTags := convertResultToTag(tagResp.Results)

	// 2ページ目以降はgoroutine
	maxPage := tagResp.Count/TAG_PER_PAGE + 1
	tagCh := make(chan types.ImageTags, maxPage-1)
	eg := errgroup.Group{}
	for page := 2; page < maxPage; page++ {
		page := page
		eg.Go(func() error {
			tagResp, err := getTagResponse(ctx, auth, option.Timeout, repository, page)
			if err != nil {
				return err
			}
			tagCh <- convertResultToTag(tagResp.Results)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for page := 2; page < maxPage; page++ {
		select {
		case tags := <-tagCh:
			imageTags = append(imageTags, tags...)
		}
	}
	return imageTags, nil
}

func convertResultToTag(summaries []ImageSummary) types.ImageTags {
	tags := []types.ImageTag{}
	for _, detail := range summaries {
		if detail.Name == "" {
			continue
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, detail.LastUpdated)
		tags = append(tags, types.ImageTag{
			Tags:      []string{detail.Name},
			Byte:      &detail.FullSize,
			CreatedAt: &createdAt,
		})
	}
	return tags
}

// getTagResponse returns the tags for a specific repository.
func getTagResponse(ctx context.Context, auth dockertypes.AuthConfig, timeout time.Duration, repository string, page int) (tagsResponse, error) {
	url := fmt.Sprintf("%s/v2/repositories/%s/tags/?page=%d", registryURL, repository, page)
	log.Logger.Debugf("url=%s,repository=%s", url, repository)
	var response tagsResponse
	if _, err := getJSON(ctx, url, auth, timeout, &response); err != nil {
		return response, err
	}

	return response, nil
}
