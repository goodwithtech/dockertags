package types

import "time"

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

type ImageTag struct {
	Tags       []string   `json:"tags"`
	Byte       *int       `json:"byte"`
	CreatedAt  *time.Time `json:"created_at"`
	UploadedAt *time.Time `json:"uploaded_at"`
}

type ImageTags []ImageTag

func (t ImageTags) Len() int      { return len(t) }
func (t ImageTags) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ImageTags) Less(i, j int) bool {
	if t[i].UploadedAt != nil && t[j].UploadedAt != nil {
		return (*(t[i].UploadedAt)).After(*(t[j].UploadedAt))
	}

	if t[i].CreatedAt != nil && t[j].CreatedAt != nil {
		return (*(t[i].CreatedAt)).After(*(t[j].CreatedAt))
	}

	return true
}
