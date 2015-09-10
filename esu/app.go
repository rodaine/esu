package esu

import (
	"io"
	"os"

	"github.com/codegangsta/cli"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "9200"
)

var (
	DefaultOutputWriter io.Writer = os.Stdout
	DefaultErrorWriter  io.Writer = os.Stderr
)

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "esu"
	app.Usage = "Elasticsearch Utility: A tool for configuring and managing an Elasticsearch cluster"
	app.Version = "0.1.0"
	app.HideHelp = true
	app.HideVersion = true
	app.Writer = DefaultOutputWriter

	setupGlobalFlags(app)
	setupGlobalCommands(app)
	setupPingCommand(app)
	setupClusterCommand(app)

	return app
}

func setupGlobalFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, h",
			Usage:  "the Elasticsearch node host",
			Value:  DefaultHost,
			EnvVar: "ESU_HOST",
		},
		cli.StringFlag{
			Name:   "port, p",
			Usage:  "the Elasticsearch node port",
			Value:  DefaultPort,
			EnvVar: "ESU_PORT",
		},
		cli.BoolFlag{
			Name:   "ssl, s",
			Usage:  "connect to the Elasticsearch node via SSL (HTTPS)",
			EnvVar: "ESU_SSL",
		},
	}
}

func setupGlobalCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "help",
			Aliases: []string{"man"},
			Usage:   "prints this help message",
			Action:  cli.ShowAppHelp,
		},
	}
}
