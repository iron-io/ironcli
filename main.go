package main

import (
	"os"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/mq"
	"github.com/iron-io/ironcli/worker"
)

func main() {
	app := cli.NewApp()
	app.Name = "iron CLI"
	app.Usage = "CLI app to access Iron.io APIs"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "version,v",
			Value:  "",
			Usage:  "the API version to use",
			EnvVar: "IRON_API_VERSION",
		},
		cli.StringFlag{
			Name:   "token,t",
			Value:  "",
			Usage:  "the OAuth token to use",
			EnvVar: "IRON_OAUTH_TOKEN",
		},
		cli.StringFlag{
			Name:   "project,p",
			Value:  "",
			Usage:  "the project ID. your OAuth token must be authorized to access the project",
			EnvVar: "IRON_PROJECT_ID",
		},
		cli.StringFlag{
			Name:   "environment,env",
			Value:  "",
			Usage:  "specify a specific dev environment",
			EnvVar: "IRON_ENVIRONMENT",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "worker",
			Aliases:     []string{"w"},
			Usage:       "IronWorker base command",
			SubCommands: worker.SubCommands,
		},
		cli.Command{
			Name:        "queue",
			Aliases:     []string{"mq"},
			Usage:       "IronMQ base command",
			SubCommands: mq.SubCommands,
		},
	}
	app.Run(os.Args)
}
