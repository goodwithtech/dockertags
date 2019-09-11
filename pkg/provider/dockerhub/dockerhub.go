package dockerhub

import (
	"context"

	"time"

	"github.com/goodwithtech/dockertags/pkg/log"
	"github.com/goodwithtech/dockertags/pkg/types"

	dockertypes "github.com/docker/cli/cli/config/types"

	"github.com/goodwithtech/dockertags/pkg/registry"
)

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
func (p *DockerHub) Run(ctx context.Context, domain, repository string, option types.AuthOption) (types.ImageTags, error) {
	auth := dockertypes.AuthConfig{
		ServerAddress: "registry.hub.docker.com",
	}
	opt := registry.Opt{}
	r, err := registry.New(ctx, auth, opt)
	if err != nil {
		return nil, err
	}
	tags, err := Tags(ctx, r, repository)
	if err != nil {
		return nil, err
	}

	imageTags := []types.ImageTag{}
	for _, detail := range tags {
		createdAt, _ := time.Parse(time.RFC3339Nano, detail.LastUpdated)
		imageTags = append(imageTags, types.ImageTag{
			Tags:      []string{detail.Name},
			Byte:      &detail.FullSize,
			CreatedAt: &createdAt,
		})
	}

	return imageTags, nil
}

// Tags returns the tags for a specific repository.
func Tags(ctx context.Context, r *registry.Registry, repository string) ([]ImageSummary, error) {
	url := r.Url("/v2/repositories/%s/tags/", repository)
	log.Logger.Debugf("url=%s,repository=%s", url, repository)
	var response tagsResponse
	if _, err := r.GetJSON(ctx, url, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}
