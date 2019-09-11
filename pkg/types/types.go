package types

import "time"

type ImageTag struct {
	Tags       []string
	Byte       *int
	CreatedAt  *time.Time
	UploadedAt *time.Time
}

var ISO8601fmt = "2006-01-02T15:04:05-0700"

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
