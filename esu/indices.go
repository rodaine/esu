package esu

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
)

func setupIndicesCommand(app *cli.App) {
	app.Commands = append(
		app.Commands,
		cli.Command{
			Name:     "indices",
			Aliases:  []string{"i", "indexes", "index"},
			Usage:    "create, read, update and delete indices",
			Action:   cli.ShowSubcommandHelp,
			HideHelp: true,
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list out info about any or all indices in the cluster",
					Action:  getIndicesList,
				},
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "create a new index",
					Action:  putIndex,
				},
				{
					Name:    "stats",
					Aliases: []string{"s"},
					Usage:   "get stats about a single index",
					Action:  getIndexStats,
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete index(es) permanently",
					Action:  deleteIndex,
				},
			},
		},
	)
}

func putIndex(ctx *cli.Context) {
	args := ctx.Args()
	var name string
	var r io.Reader

	switch len(args) {
	case 1:
		name = args[0]
		r = getStdIn()
	case 2:
		name = args[0]
		r = getFile(args[1])
	default:
		cli.ShowSubcommandHelp(ctx)
		os.Exit(1)
	}

	svc := connectToES(ctx).CreateIndex(name)

	if r != nil {
		settings, err := readJSON(r)
		if err != nil {
			exitWithError(err)
		}
		svc.BodyJson(settings)
	}

	res, err := svc.Do()
	if err != nil {
		exitWithError(err)
	}

	t := NewTable("Index Creation", name)
	t.Add("Acknowledged", res.Acknowledged)
	t.Print(ctx)
}

func getIndicesList(ctx *cli.Context) {
	args := []string(ctx.Args())
	if len(args) == 0 {
		args = []string{"_all"}
	}

	es := connectToES(ctx)

	settings, err := es.IndexGetSettings(args...).FlatSettings(true).Do()
	if err != nil {
		exitWithError(err)
	}

	stats, err := connectToES(ctx).IndexStats(args...).Metric("docs").Do()
	if err != nil {
		exitWithError(err)
	}

	t := NewTable("Name", "Shards", "Replicas", "Documents")
	for name, info := range stats.Indices {
		setting := settings[name].Settings
		t.Add(
			name,
			setting["index.number_of_shards"],
			setting["index.number_of_replicas"],
			fmt.Sprintf("%d (%d Deleted)", info.Primaries.Docs.Count, info.Primaries.Docs.Deleted),
		)
	}

	t.Print(ctx)
}

func getIndexStats(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 1 {
		exitWithHelp(ctx)
	}
	idx := args.First()

	res, err := connectToES(ctx).IndexStats(idx).Human(true).Do()
	if err != nil {
		exitWithError(err)
	}

	stats, ok := res.Indices[idx]
	if !ok {
		exitWithError(fmt.Errorf("Unable to find index: %s", idx))
	}

	var t *Table

	if stats.Primaries.Docs != nil {
		t = NewTable("Documents", "")
		t.Add("Total", stats.Primaries.Docs.Count)
		t.Add("Deleted", stats.Primaries.Docs.Deleted)
		t.Print(ctx)
	}

	if stats.Primaries.Store != nil {
		t = NewTable("Storage", "")
		t.Add("Size", stats.Primaries.Store.Size)
		t.Add("Throttle Time", stats.Primaries.Store.ThrottleTime)
		t.Print(ctx)
	}
}

func deleteIndex(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) == 0 {
		exitWithHelp(ctx)
	}
}
