package dockerhub

import (
	"context"
	"fmt"
	"sort"

	"golang.org/x/sync/errgroup"

	"time"

	"github.com/goodwithtech/dockertags/internal/log"
	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"

	dockertypes "github.com/docker/docker/api/types"
)

const registryURL = "https://registry.hub.docker.com/"

// DockerHub implements Run
type DockerHub struct {
	filterOpt *types.FilterOption
}

// Run returns tag list
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

	// create all in one []ImageSummary
	totalTagSummary := tagResp.Results
	lastPage := calcMaxRequestPage(tagResp.Count, reqOpt.MaxCount, filterOpt)
	// create ch (page - 1), already fetched first page,
	tagsPerPage := make(chan []ImageSummary, lastPage-1)
	eg := errgroup.Group{}
	for page := 2; page <= lastPage; page++ {
		page := page
		eg.Go(func() error {
			tagResp, err := getTagResponse(ctx, auth, reqOpt.Timeout, repository, page)
			if err != nil {
				return err
			}
			tagsPerPage <- tagResp.Results
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for page := 2; page <= lastPage; page++ {
		select {
		case tags := <-tagsPerPage:
			totalTagSummary = append(totalTagSummary, tags...)
		}
	}
	return p.convertResultToTag(totalTagSummary), nil
}

func summarizeByHash(summaries []ImageSummary) map[string]types.ImageTag {
	pools := map[string]types.ImageTag{}
	for _, imageSummary := range summaries {
		if imageSummary.Name == "" {
			log.Logger.Debugf("no tag data :%v", imageSummary)
			continue
		}
		if len(imageSummary.Images) == 0 {
			log.Logger.Debugf("no image layer data :%v", imageSummary)
			continue
		}
		sort.Sort(imageSummary.Images)
		firstHash := imageSummary.Images[0].Digest
		target, ok := pools[firstHash]
		// create first one if not exist
		if !ok {
			pools[firstHash] = createImageTag(imageSummary)
			continue
		}
		// set newer CreatedAt
		target.Tags = append(target.Tags, imageSummary.Name)
		createdAt, _ := time.Parse(time.RFC3339Nano, imageSummary.LastUpdated)
		if createdAt.After(target.CreatedAt) {
			target.CreatedAt = createdAt
		}
		pools[firstHash] = target
	}
	return pools
}

func (p *DockerHub) convertResultToTag(summaries []ImageSummary) types.ImageTags {
	// create map : key is image hash
	pools := summarizeByHash(summaries)
	tags := []types.ImageTag{}
	for _, imageTag := range pools {
		if !utils.MatchConditionTags(p.filterOpt, imageTag.Tags) {
			continue
		}
		tags = append(tags, imageTag)
	}
	return tags
}

func createImageTag(is ImageSummary) types.ImageTag {
	createdAt, _ := time.Parse(time.RFC3339Nano, is.LastUpdated)
	tagNames := []string{is.Name}
	return types.ImageTag{
		Tags:      tagNames,
		Byte:      is.FullSize,
		CreatedAt: createdAt,
	}
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
	return totalCnt/types.ImagePerPage + 1
	// TODO : currently always fetch all pages for show alias
	// maxPage := totalCnt/types.ImagePerPage + 1
	// if needCnt == 0 || len(option.Contain) != 0 {
	// 	return maxPage
	// }
	// needPage := needCnt/types.ImagePerPage + 1
	// if needPage >= maxPage {
	// 	return maxPage
	// }
	// return needPage
}
