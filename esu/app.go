package esu

import "github.com/codegangsta/cli"

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "esu"
	app.Usage = "ElasticSearch Utility: A tool for configuring and managing an ElasticSearch cluster"
	app.Version = "0.1.0"
	return app
}
