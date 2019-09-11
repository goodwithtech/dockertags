package gcr

import (
	"context"
	"strconv"
	"time"

	"github.com/goodwithtech/dockertags/pkg/log"

	"github.com/goodwithtech/dockertags/pkg/types"

	"github.com/GoogleCloudPlatform/docker-credential-gcr/config"
	"github.com/GoogleCloudPlatform/docker-credential-gcr/credhelper"
	"github.com/GoogleCloudPlatform/docker-credential-gcr/store"
	dockertypes "github.com/docker/cli/cli/config/types"
	"github.com/goodwithtech/dockertags/pkg/registry"
)

type GCR struct {
	registry *registry.Registry
	domain   string
	Store    store.GCRCredStore
}

type tagsResponse struct {
	Next     string                     `json:"next"`
	Previous string                     `json:"previous"`
	Manifest map[string]ManifestSummary `json:"manifest"`
}

type ManifestSummary struct {
	Tag            []string `json:"tag"`
	ImageSizeBytes string   `json:"imageSizeBytes"`
	CreatedMS      string   `json:"timeCreatedMs"`
	UploadedMS     string   `json:"timeUploadedMs"`
}

func (p *GCR) Run(ctx context.Context, domain, repository string, option types.AuthOption) (imageTags types.ImageTags, err error) {
	if option.GcpCredPath != "" {
		p.Store = store.NewGCRCredStore(option.GcpCredPath)
	}
	auth := dockertypes.AuthConfig{
		ServerAddress: domain,
	}
	auth.Username, auth.Password, err = p.getCredential(ctx)
	if err != nil {
		log.Logger.Debug("Fail to get gcp credential")
		//return nil, err
	}

	opt := registry.Opt{Timeout: option.Timeout}
	r, err := registry.New(ctx, auth, opt)
	if err != nil {
		return nil, err
	}
	p.registry = r
	tags, err := p.getTags(ctx, repository)
	if err != nil {
		return nil, err
	}

	for _, detail := range tags {
		if len(detail.Tag) == 0 {
			continue
		}
		createdAt, err := stringMStoTime(detail.CreatedMS)
		if err != nil {
			return nil, err
		}

		uploadedAt, err := stringMStoTime(detail.UploadedMS)
		if err != nil {
			return nil, err
		}

		imageTags = append(imageTags, types.ImageTag{
			Tags:       detail.Tag,
			Byte:       getIntByte(detail.ImageSizeBytes),
			CreatedAt:  &createdAt,
			UploadedAt: &uploadedAt,
		})
	}
	return imageTags, nil
}

func stringMStoTime(msstring string) (time.Time, error) {
	intMS, err := strconv.Atoi(msstring)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := int64(intMS / 1000)
	return time.Unix(timestamp, 0), nil
}

// getTags returns the tags
func (p *GCR) getTags(ctx context.Context, repository string) (map[string]ManifestSummary, error) {
	url := p.registry.Url("/v2/%s/tags/list", repository)
	log.Logger.Debugf("url=%s,repository=%s", url, repository)

	var response tagsResponse
	if _, err := p.registry.GetJSON(ctx, url, &response); err != nil {
		return nil, err
	}

	return response.Manifest, nil
}

func (g *GCR) getCredential(ctx context.Context) (username, password string, err error) {
	var credStore store.GCRCredStore
	if g.Store == nil {
		credStore, err = store.DefaultGCRCredStore()
		if err != nil {
			return "", "", err
		}
	} else {
		credStore = g.Store
	}
	userCfg, err := config.LoadUserConfig()
	if err != nil {
		return "", "", err
	}
	helper := credhelper.NewGCRCredentialHelper(credStore, userCfg)
	return helper.Get(g.domain)
}

func getIntByte(strB string) *int {
	b, err := strconv.Atoi(strB)
	if err != nil {
		return nil
	}
	return &b
}
