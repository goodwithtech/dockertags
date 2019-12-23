package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/goodwithtech/dockertags/internal/types"
)

type JsonWriter struct {
	Output             io.Writer
	RepositoryName     string
	ShowRepositoryName bool
}

type OutputJSON struct {
	Repository string     `json:"tags,"`
	Tags       []string   `json:"tags"`
	Byte       *int       `json:"byte"`
	CreatedAt  *time.Time `json:"created_at"`
	UploadedAt *time.Time `json:"uploaded_at"`
}

func (jw JsonWriter) Write(tags types.ImageTags) (err error) {
	output, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if _, err = fmt.Fprint(jw.Output, string(output)); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}
	return nil
}
