package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodwithtech/dockertags/pkg/log"
	"github.com/goodwithtech/dockertags/pkg/provider"
	"github.com/goodwithtech/dockertags/pkg/types"
	"github.com/goodwithtech/dockertags/pkg/utils"
	"github.com/olekukonko/tablewriter"

	l "log"
	"time"

	"github.com/urfave/cli"
)

func Run(c *cli.Context) (err error) {
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	cliVersion := "v" + c.App.Version
	latestVersion, err := utils.FetchLatestVersion()

	// check latest version
	if err == nil && cliVersion != latestVersion && c.App.Version != "dev" {
		log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
	}
	err = nil
	args := c.Args()
	image := args[0]
	if image == "" {
		return fmt.Errorf("input any image name")
	}
	opt := types.AuthOption{
		Timeout: time.Second * 10,
	}

	log.Logger.Debug("Start fetch tags...")
	tags, err := provider.Exec(image, opt)
	if err != nil {
		return err
	}

	log.Logger.Debug("Writing table...")
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

	return nil
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
	return utils.ByteSize(*b)
}

func ttos(t *time.Time) string {
	if t == nil {
		return "NULL"
	}
	return (*t).Format(time.RFC3339)
}
