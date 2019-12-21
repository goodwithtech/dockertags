package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goodwithtech/dockertags/internal/types"
)

type JsonWriter struct {
	Output io.Writer
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
