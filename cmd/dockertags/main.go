package main

import (
	l "log"
	"os"
	"time"

	"github.com/goodwithtech/dockertags/pkg"
	"github.com/goodwithtech/dockertags/pkg/log"
	"github.com/urfave/cli"
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

	app.Usage = "Fetch docker tags informations"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "fetch all tagged image information",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Value: 50,
			Usage: "Set max fetch count",
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
		cli.StringFlag{
			Name:  "authurl, auth",
			Usage: "Url when fetch authentication",
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
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
