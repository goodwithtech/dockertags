package dockerhub

// dockerhub api's format
type tagsResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []ImageSummary `json:"results"`
}

// ImageSummary depends on dockerhub api
type ImageSummary struct {
	Name        string `json:"name"`
	FullSize    int    `json:"full_size"`
	LastUpdated string `json:"last_updated"`
	Images      Images `json:"images"`
}

// Images is slice of Image
type Images []Image

// Image implement sort by digest hash for detect same images
type Image struct {
	Digest       string `json:"digest"`
	Architecture string `json:"architecture"`
}

// Len is
func (t Images) Len() int { return len(t) }

// Swap is
func (t Images) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less : image sort by digest hash
func (t Images) Less(i, j int) bool {
	return (t[i].Digest) > (t[j].Digest)
}
