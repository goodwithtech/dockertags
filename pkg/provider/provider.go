package provider

import (
	"context"
	"sort"
	"strings"

	"github.com/goodwithtech/dockertags/pkg/image"

	"github.com/goodwithtech/dockertags/pkg/types"

	"github.com/goodwithtech/dockertags/pkg/provider/dockerhub"
	"github.com/goodwithtech/dockertags/pkg/provider/ecr"
	"github.com/goodwithtech/dockertags/pkg/provider/gcr"
)

const (
	ecrURL = "amazonaws.com"
	gcrURL = "gcr.io"
)

type Provider interface {
	Run(ctx context.Context, domain, repository string, option types.AuthOption) (types.ImageTags, error)
}

func Exec(imageName string, option types.AuthOption) (types.ImageTags, error) {
	image, err := image.ParseImage(imageName)
	if err != nil {
		return nil, err
	}

	p := NewProvider(image.Domain)
	ctx, cancel := context.WithTimeout(context.Background(), option.Timeout)
	defer cancel()
	imageTags, err := p.Run(ctx, image.Domain, image.Path, option)
	if err != nil {
		return nil, err
	}
	sort.Sort(types.ImageTags(imageTags))
	return imageTags, nil
}

func NewProvider(domain string) Provider {
	if strings.HasSuffix(domain, ecrURL) {
		return &ecr.ECR{}
	}
	if strings.HasSuffix(domain, gcrURL) {
		return &gcr.GCR{}
	}
	return &dockerhub.DockerHub{}
}
