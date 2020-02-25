package report

import "github.com/goodwithtech/dockertags/internal/types"

// Writer is
type Writer interface {
	Write(types.ImageTags) error
}
