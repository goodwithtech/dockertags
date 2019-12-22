package provider

import (
	"context"
	"sort"
	"strings"

	"github.com/goodwithtech/dockertags/pkg/image"

	"github.com/goodwithtech/dockertags/internal/types"

	"github.com/goodwithtech/dockertags/pkg/provider/dockerhub"
	"github.com/goodwithtech/dockertags/pkg/provider/ecr"
	"github.com/goodwithtech/dockertags/pkg/provider/gcr"
)

var containNeedle string

const (
	ecrURL = "amazonaws.com"
	gcrURL = "gcr.io"
)

type Provider interface {
	Run(ctx context.Context, domain, repository string, reqOpt *types.RequestOption, filterOpt *types.FilterOption) (types.ImageTags, error)
}

func Exec(imageName string, reqOpt *types.RequestOption, filterOpt *types.FilterOption) (types.ImageTags, error) {
	image, err := image.ParseImage(imageName)
	if err != nil {
		return nil, err
	}

	p := NewProvider(image.Domain)
	ctx, cancel := context.WithTimeout(context.Background(), reqOpt.Timeout)
	defer cancel()
	imageTags, err := p.Run(ctx, image.Domain, image.Path, reqOpt, filterOpt)
	if err != nil {
		return nil, err
	}
	sort.Sort(imageTags)
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
