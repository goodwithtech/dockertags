package report

import "github.com/goodwithtech/dockertags/internal/types"

type Writer interface {
	Write(types.ImageTags) error
}
