package main

import (
	l "log"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/goodwithtech/dockertags/internal/log"
	"github.com/goodwithtech/dockertags/pkg"
)

var (
	version = "dev"
)

func main() {
	cli.AppHelpTemplate = `NAME:
  {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
  {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
VERSION:
  {{.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
  {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
  {{range $index, $author := .Authors}}{{if $index}}
  {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
OPTIONS:
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}
`
	app := cli.NewApp()
	app.Name = "dockertags"
	app.Version = version
	app.ArgsUsage = "image_name"

	app.Usage = "fetch docker image tags"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "set max tags count. if exist no tag image will be short numbers. limit=0 means fetch all tags",
		},
		cli.StringSliceFlag{
			Name:  "contain, c",
			Usage: "contains target string. multiple string allows.",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "table",
			Usage: "target format table or json, default table",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output file name, default output to stdout",
		},
		cli.StringFlag{
			Name:  "authurl, auth",
			Usage: "GetURL when fetch authentication",
		},
		cli.DurationFlag{
			Name:  "timeout, t",
			Value: time.Second * 10,
			Usage: "e.g)5s, 1m",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "Username",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "Using -password via CLI is insecure. Be careful.",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Show debug logs",
		},
	}

	app.Action = pkg.Run
	err := app.Run(os.Args)
	if err != nil {
		if log.Logger != nil {
			log.Logger.Fatal(err)
		}
		l.Fatal(err)
	}
}
