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

const (
	registryURL = "https://registry.hub.docker.com/"
	rateLimit   = 64
)

// DockerHub implements Run
type DockerHub struct {
	filterOpt  *types.FilterOption
	requestOpt *types.RequestOption
	authCfg    dockertypes.AuthConfig
}

// Run returns tag list
func (p *DockerHub) Run(ctx context.Context, domain, repository string, reqOpt *types.RequestOption, filterOpt *types.FilterOption) (types.ImageTags, error) {
	p.filterOpt = filterOpt
	p.requestOpt = reqOpt
	p.authCfg = dockertypes.AuthConfig{
		ServerAddress: "registry.hub.docker.com",
		Username:      reqOpt.UserName,
		Password:      reqOpt.Password,
	}

	// fetch page 1 for check max item count.
	tagResp, err := getTagResponse(ctx, p.authCfg, reqOpt.Timeout, repository, 1)
	if err != nil {
		return nil, err
	}

	// create all in one []tagSummary
	totalTagSummary := tagResp.Results
	lastPage := calcMaxRequestPage(tagResp.Count, reqOpt.MaxCount, filterOpt)
	// create ch (page - 1), already fetched first page,
	tagsPerPage := make(chan []tagSummary, lastPage-1)
	if err = p.controlGetTags(ctx, tagsPerPage, repository, 2, lastPage); err != nil {
		return nil, err
	}
	for page := 2; page <= lastPage; page++ {
		select {
		case tags := <-tagsPerPage:
			totalTagSummary = append(totalTagSummary, tags...)
		}
	}
	close(tagsPerPage)
	return p.convertResultToTag(totalTagSummary), nil
}

// rate limit for socket: too many open files
func (p *DockerHub) controlGetTags(ctx context.Context, tagsPerPage chan []tagSummary, repository string, from, to int) error {
	slots := make(chan struct{}, rateLimit)
	eg := errgroup.Group{}
	for page := from; page <= to; page++ {
		page := page
		slots <- struct{}{}
		eg.Go(func() error {
			defer func() { <-slots }()
			tagResp, err := getTagResponse(ctx, p.authCfg, p.requestOpt.Timeout, repository, page)
			if err != nil {
				return err
			}
			tagsPerPage <- tagResp.Results
			return nil
		})
	}
	return eg.Wait()
}

func summarizeByHash(summaries []tagSummary) map[string]types.ImageTag {
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
		for _, img := range imageSummary.Images {
			hash := img.Digest
			target, ok := pools[hash]
			// create first hash key if not exist
			if !ok {
				pools[hash] = convertUploadImageTag(imageSummary, img)
				continue
			}
			// set newer uploaded at
			target.Tags = append(target.Tags, imageSummary.Name)
			target.OsArchs = append(target.OsArchs, genOsArch(img))
			uploadedAt, _ := time.Parse(time.RFC3339Nano, imageSummary.LastUpdated)
			if uploadedAt.After(target.UploadedAt) {
				target.UploadedAt = uploadedAt
			}
			pools[hash] = target
		}
	}
	return pools
}

func genOsArch(img image) string {
	return fmt.Sprintf("%s/%s", img.Os, img.Architecture)
}

func (p *DockerHub) convertResultToTag(summaries []tagSummary) types.ImageTags {
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

func convertUploadImageTag(is tagSummary, img image) types.ImageTag {
	uploadedAt, _ := time.Parse(time.RFC3339Nano, is.LastUpdated)
	tagNames := []string{is.Name}
	return types.ImageTag{
		Tags:       tagNames,
		Byte:       is.FullSize,
		Hash:       img.Digest,
		OsArchs:    []string{genOsArch(img)},
		UploadedAt: uploadedAt,
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
	//return totalCnt/types.ImagePerPage + 1
	maxPage := totalCnt / types.ImagePerPage
	if totalCnt%types.ImagePerPage != 0 {
		maxPage++
	}
	if needCnt == 0 || len(option.Contain) != 0 {
		return maxPage
	}

	needPage := needCnt / types.ImagePerPage
	if needCnt%types.ImagePerPage != 0 {
		needCnt++
	}
	if needPage >= maxPage {
		return maxPage
	}
	return needPage
}
