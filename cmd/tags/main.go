package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goodwithtech/image-tag-sorter/pkg/util"

	"github.com/goodwithtech/image-tag-sorter/pkg/types"
	"github.com/olekukonko/tablewriter"

	"github.com/goodwithtech/image-tag-sorter/pkg/provider"
)

func main() {

	image := "goodwithtech/dockle"
	opt := types.AuthOption{
		Timeout: time.Second * 10,
	}
	tags, err := provider.Exec(image, opt)
	if err != nil {
		fmt.Println("err", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Full", "Tag", "Size", "Created At", "Uploaded At"})

	for _, tag := range tags {
		table.Append([]string{
			getFullPath(image, tag.Tags),
			strings.Join(tag.Tags, ","),
			getBytesize(tag.Byte),
			ttos(tag.CreatedAt),
			ttos(tag.UploadedAt),
		})
	}
	table.Render()

}
func getFullPath(image string, tags []string) string {
	if len(tags) == 0 {
		return "NO TAGGED"
	}
	return image + ":" + tags[0]
}

func getBytesize(b *int) string {
	if b == nil {
		return "-"
	}
	return util.ByteSize(*b)
}

func ttos(t *time.Time) string {
	if t == nil {
		return "NULL"
	}
	return (*t).Format(time.RFC3339)
}
