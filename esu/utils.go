package esu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"gopkg.in/olivere/elastic.v3-unstable"
)

func getConnectionURL(ctx *cli.Context) *url.URL {
	scheme := "http"
	if ctx.GlobalBool("ssl") {
		scheme = "https"
	}

	return &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%s", ctx.GlobalString("host"), ctx.GlobalString("port")),
	}
}

func connectToES(ctx *cli.Context) (es *elastic.Client) {
	uri := getConnectionURL(ctx)
	es, err := elastic.NewClient(
		elastic.SetURL(uri.String()),
		elastic.SetSniff(false),       // esu might not be able to ping other nodes
		elastic.SetHealthcheck(false), // healthchecks are a skosh overkill for a CLI tool
	)

	if err != nil {
		exitWithError(err)
	}

	return
}

func readJSON(r io.Reader) (out map[string]interface{}, err error) {
	d := json.NewDecoder(r)
	err = d.Decode(&out)
	return
}

func exitWithError(err error) {
	txt := color.New(color.FgRed).SprintfFunc()("ERROR: %v", err)
	fmt.Fprintln(DefaultErrorWriter, txt)
	os.Exit(1)
}
