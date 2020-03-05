package types

import "time"

// ImagePerPage :
const ImagePerPage = 10

// RequestOption : container registry information
type RequestOption struct {
	MaxCount        int
	AuthURL         string
	UserName        string
	Password        string
	GcpCredPath     string
	AwsAccessKey    string
	AwsSecretKey    string
	AwsSessionToken string
	AwsRegion       string
	Timeout         time.Duration
}

// FilterOption : tag pattern
type FilterOption struct {
	Contain []string
}

// ImageTag : tag information
type ImageTag struct {
	Tags       []string  `json:"tags"`
	Data       []TagAttr `json:"data"`
	CreatedAt  time.Time `json:"created_at"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// TagAttr is each unique tag information
type TagAttr struct {
	Os     string
	Arch   string
	Digest string
	Byte   int `json:"byte"`
}

// ImageTags : tag information slice
type ImageTags []ImageTag

// Len interface method of sort
func (t ImageTags) Len() int      { return len(t) }
func (t ImageTags) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ImageTags) Less(i, j int) bool {
	if !t[i].UploadedAt.IsZero() && !t[j].UploadedAt.IsZero() {
		return (t[i].UploadedAt).After((t[j].UploadedAt))
	}
	if !t[i].CreatedAt.IsZero() && !t[j].CreatedAt.IsZero() {
		return (t[i].CreatedAt).After((t[j].CreatedAt))
	}
	return true
}

// TagAttrs is each unique tag information
type TagAttrs []TagAttr

func (t TagAttrs) Len() int           { return len(t) }
func (t TagAttrs) Less(i, j int) bool { return t[i].Digest < t[j].Digest }
func (t TagAttrs) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
