package report

import (
	"fmt"
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
	Output      io.Writer
	LongDigests bool
}

// Write is
func (w TableWriter) Write(tags types.ImageTags) (err error) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tag", "Size", "Digest", "OS/ARCH", "Created At", "Uploaded At"})

	for _, tag := range tags {
		targets := utils.StrByLen(tag.Tags)
		sort.Sort(targets)

		var sizes, digests, osArchs []string
		for _, datum := range tag.Data {
			sizes = append(sizes, getBytesize(datum.Byte))
			digest := datum.Digest
			if !w.LongDigests {
				digest = trimHash(datum.Digest)
			}
			digests = append(digests, digest)
			if datum.Os != "" {
				osArchs = append(osArchs, fmt.Sprintf("%s/%s", datum.Os, datum.Arch))
			}
		}
		// filled with whitespace
		table.Append([]string{
			strings.Join(fillWithSpaces(targets), tablewriter.NEWLINE),
			strings.Join(fillWithSpaces(sizes), tablewriter.NEWLINE),
			strings.Join(fillWithSpaces(digests), tablewriter.NEWLINE),
			strings.Join(fillWithSpaces(osArchs), tablewriter.NEWLINE),
			ttos(tag.CreatedAt),
			ttos(tag.UploadedAt),
		})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.Render()
	return nil
}

func fillWithSpaces(labels []string) []string {
	fillWithSpaces := []string{}
	for _, tag := range labels {
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
