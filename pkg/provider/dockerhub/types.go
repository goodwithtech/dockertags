package dockerhub

// dockerhub api's format
type tagsResponse struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []tagSummary `json:"results"`
}

// tagSummary depends on dockerhub api
type tagSummary struct {
	Name        string `json:"name"`
	FullSize    int    `json:"full_size"`
	LastUpdated string `json:"last_updated"`
	Images      images `json:"images"`
}

// images is slice of image
type images []image

// image implement sort by digest hash for detect same images
type image struct {
	Digest       string `json:"digest"`
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
}

// Len is
func (t images) Len() int { return len(t) }

// Swap is
func (t images) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less : image sort by digest hash
func (t images) Less(i, j int) bool {
	return (t[i].Digest) > (t[j].Digest)
}
