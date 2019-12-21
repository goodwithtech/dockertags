package report

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"
)

type TableWriter struct {
	Output io.Writer
}

func (w TableWriter) Write(tags types.ImageTags) (err error) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tag", "Size", "Created At", "Uploaded At"})
	for _, tag := range tags {
		table.Append([]string{
			strings.Join(tag.Tags, ","),
			getBytesize(tag.Byte),
			ttos(tag.CreatedAt),
			ttos(tag.UploadedAt),
		})
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