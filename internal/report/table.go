package report

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"
)

type TableWriter struct {
	Output             io.Writer
	RepositoryName     string
	ShowRepositoryName bool
}

func (w TableWriter) Write(tags types.ImageTags) (err error) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Tags", "Size", "Created At", "Uploaded At"})
	if w.ShowRepositoryName {
		table.SetHeader([]string{"Repository", "Tags", "Size", "Created At", "Uploaded At"})
	}

	for _, tag := range tags {
		content := []string{
			strings.Join(tag.Tags, ","),
			getBytesize(tag.Byte),
			ttos(tag.CreatedAt),
			ttos(tag.UploadedAt),
		}
		if w.ShowRepositoryName {
			if len(tag.Tags) == 0 {
				continue
			}
			content = append([]string{fmt.Sprintf("%s:%s", w.RepositoryName, tag.Tags[0])}, content...)
		}
		table.Append(content)
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return nil
}

func getBytesize(b *int) string {
	if b == nil {
		return "-"
	}
	return utils.ByteSize(*b)
}

func ttos(t *time.Time) string {
	if t == nil {
		return "NULL"
	}
	return (*t).Format(time.RFC3339)
}
