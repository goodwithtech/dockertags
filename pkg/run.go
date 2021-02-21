package pkg

import (
	"context"
	"errors"
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

// Run runs dockertag operation
func Run(c *cli.Context) (err error) {
	debug := c.Bool("debug")
	// reload logger if set flag
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}
	cliVersion := "v" + c.App.Version
	ctx, cancel := context.WithTimeout(context.Background(), c.Duration("timeout"))
	defer cancel()
	latestVersion, err := utils.FetchLatestVersion(ctx)
	// check latest version
	if err != nil {
		log.Logger.Debugf("Failed to check latest version. %s", err)
	} else if cliVersion != latestVersion && c.App.Version != "dev" {
		log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
	}

	var imageName string
	if imageName, err = fetchImageName(c.Args()); err != nil {
		log.Logger.Fatalf("%s", err)
	}

	reqOpt, filterOpt, err := genOpts(c)
	if err != nil {
		log.Logger.Fatalf("invalid option: %s.", err)
	}

	tags, err := provider.Exec(imageName, reqOpt, filterOpt)
	if err != nil {
		return err
	}
	log.Logger.Debug("Start reporting")

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
		writer = &report.JSONWriter{Output: output}
	default:
		longDigests := c.Bool("digests")
		writer = &report.TableWriter{Output: output, LongDigests: longDigests}
	}

	var showTags types.ImageTags
	if reqOpt.MaxCount > 0 && reqOpt.MaxCount < len(tags) {
		showTags = tags[:reqOpt.MaxCount]
	} else {
		showTags = tags
	}
	return writer.Write(showTags)
}

func genRequestOpt(c *cli.Context) types.RequestOption {
	return types.RequestOption{
		MaxCount: c.Int("limit"),
		Timeout:  c.Duration("timeout"),
		AuthURL:  c.String("authurl"),
		UserName: c.String("username"),
		Password: c.String("password"),
	}
}

func genFilterOpt(c *cli.Context) types.FilterOption {
	return types.FilterOption{
		Contain: c.StringSlice("contain"),
	}
}

func genOpts(c *cli.Context) (*types.RequestOption, *types.FilterOption, error) {
	reqOpt := genRequestOpt(c)
	filterOpt := genFilterOpt(c)
	return &reqOpt, &filterOpt, nil
}

func fetchImageName(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("dockertags requires one argument")
	}
	return args[0], nil
}
