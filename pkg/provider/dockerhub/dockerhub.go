package dockerhub

import (
	"context"
	"fmt"
	"sort"
	"strings"

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
			uploadedAt, _ := time.Parse(time.RFC3339Nano, imageSummary.LastUpdated)
			if uploadedAt.After(target.UploadedAt) {
				target.UploadedAt = uploadedAt
			}
			pools[hash] = target
		}
	}
	return pools
}

func (p *DockerHub) convertResultToTag(summaries []tagSummary) types.ImageTags {
	// create map : key is image hash
	keyDigestMap := summarizeByHash(summaries)

	// latestv3.3.9 => []ImageTag
	keyTagsMap := p.summarizeByTagNames(keyDigestMap)

	// []ImageTags
	imageTags := make([]types.ImageTag, 0, len(keyTagsMap))
	for _, digests := range keyTagsMap {
		imageTag := types.ImageTag{}
		for idx, digest := range digests {
			targetImageTag := keyDigestMap[digest]
			if idx == 0 {
				imageTag = targetImageTag
				continue
			}
			imageTag.Data = append(imageTag.Data, targetImageTag.Data[0])
		}
		imageTags = append(imageTags, imageTag)
	}
	return imageTags
}

func (p *DockerHub) summarizeByTagNames(keyDigestMap map[string]types.ImageTag) map[string][]string {
	keyTagsMap := map[string][]string{}
	for digest, imageTag := range keyDigestMap {
		if !utils.MatchConditionTags(p.filterOpt, imageTag.Tags) {
			continue
		}
		key := strings.Join(imageTag.Tags, "")
		digests, ok := keyTagsMap[key]
		if !ok {
			keyTagsMap[key] = []string{digest}
		}
		keyTagsMap[key] = append(digests, digest)
	}
	return keyTagsMap
}

func convertUploadImageTag(is tagSummary, img image) types.ImageTag {
	uploadedAt, _ := time.Parse(time.RFC3339Nano, is.LastUpdated)
	tagNames := []string{is.Name}
	return types.ImageTag{
		Tags: tagNames,
		Data: []types.TagAttr{{
			Os:     img.Os,
			Arch:   img.Architecture,
			Digest: img.Digest,
			Byte:   img.Size,
		}},
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
		// only notice error
		log.Logger.Errorf("invalid response at %s", url)
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
