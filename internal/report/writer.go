package report

import "github.com/goodwithtech/dockertags/internal/types"

// Writer is
type Writer interface {
	Write(types.ImageTags) error
}

func trimHash(long string) string {
	if len(long) < 18 {
		return long
	}
	if long[0:6] == "sha256" {
		return long[7:16]
	}
	return long
}
