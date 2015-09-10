package esu

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func setupPingCommand(app *cli.App) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "ping the cluster to see if it's available",
		Action:  ping,
	},
	)
}

func ping(ctx *cli.Context) {
	uri := getConnectionURL(ctx)
	res, _, err := connectToES(ctx).Ping(uri.String()).Do()

	if err != nil {
		exitWithError(err)
	}

	t := NewTable("Cluster", res.ClusterName)
	t.Add("Node", fmt.Sprintf("%s [%v]", res.Name, uri))
	t.Add("Tag Line", res.TagLine)
	t.Add("ES Version", res.Version.Number)
	t.Print(ctx)
}
