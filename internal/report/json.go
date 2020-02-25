package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goodwithtech/dockertags/internal/types"
)

// JSONWriter create json output
type JSONWriter struct {
	Output io.Writer
}

// Write is
func (jw JSONWriter) Write(tags types.ImageTags) (err error) {
	output, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if _, err = fmt.Fprint(jw.Output, string(output)); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}
	return nil
}
