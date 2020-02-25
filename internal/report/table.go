package report

import (
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"
)

// TableWriter output table format
type TableWriter struct {
	Output io.Writer
}

// Write is
func (w TableWriter) Write(tags types.ImageTags) (err error) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tag", "Size", "Created At", "Uploaded At"})

	for _, tag := range tags {
		targets := utils.StrByLen(tag.Tags)
		sort.Sort(targets)

		// filled with whitespace
		table.Append([]string{
			strings.Join(fillWithSpaces(targets), tablewriter.NEWLINE),
			getBytesize(tag.Byte),
			ttos(tag.CreatedAt),
			ttos(tag.UploadedAt),
		})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.Render()

	return nil
}

func fillWithSpaces(tagNames []string) []string {
	fillWithSpaces := []string{}
	for _, tag := range tagNames {
		tagStr := tag
		whitespaceCnt := tablewriter.MAX_ROW_WIDTH - len(tag)
		if whitespaceCnt > 0 {
			tagStr = tag + strings.Repeat(tablewriter.SPACE, whitespaceCnt)
		}
		fillWithSpaces = append(fillWithSpaces, tagStr)
	}
	return fillWithSpaces
}

func getBytesize(b int) string {
	if b == 0 {
		return "-"
	}
	return utils.ByteSize(b)
}

func ttos(t time.Time) string {
	if t.IsZero() {
		return "NULL"
	}
	return t.Format(time.RFC3339)
}
