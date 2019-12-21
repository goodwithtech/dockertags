package pkg

import (
	"fmt"
	l "log"
	"os"

	"github.com/goodwithtech/dockertags/internal/log"
	"github.com/goodwithtech/dockertags/internal/report"
	"github.com/goodwithtech/dockertags/internal/types"
	"github.com/goodwithtech/dockertags/internal/utils"
	"github.com/goodwithtech/dockertags/pkg/provider"

	"github.com/urfave/cli"
)

func Run(c *cli.Context) (err error) {
	debug := c.Bool("debug")
	all := c.Bool("all")
	// reload logger if set flag
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

	if len(args) == 0 {
		log.Logger.Infof(`"dockertags" requires at least 1 argument.`)
		cli.ShowAppHelpAndExit(c, 1)
		return
	}
	image := args[0]
	opt := types.RequestOption{
		MaxCount: c.Int("limit"),
		Timeout:  c.Duration("timeout"),
		AuthURL:  c.String("authurl"),
		UserName: c.String("username"),
		Password: c.String("password"),
	}

	log.Logger.Debug("Start fetch tags...")
	tags, err := provider.Exec(image, opt)
	if err != nil {
		return err
	}
	var showTags types.ImageTags
	showTags = tags
	if !all && opt.MaxCount < len(tags) {
		showTags = tags[:opt.MaxCount]
	}

	log.Logger.Debug("Writing table...")

	o := c.String("output")
	output := os.Stdout
	if o != "" {
		if output, err = os.Create(o); err != nil {
			return fmt.Errorf("failed to create an output file: %w", err)
		}
	}
	var writer report.Writer
	switch format := c.String("format"); format {
	case "json":
		writer = &report.JsonWriter{Output: output}
	default:
		writer = &report.TableWriter{Output: output}
	}

	return writer.Write(showTags)
}
