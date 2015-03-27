package main

import (
	"fmt"
	"os"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/ironcli/mq"
	"github.com/iron-io/ironcli/worker"
)

func main() {
	app := cli.NewApp()
	app.Name = "ironcli"
	app.Usage = "CLI app to access Iron.io APIs"
	app.Version = common.VersionNum
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   common.Version,
			Usage:  "the API version to use",
			Value:  common.InvalidVersion,
			EnvVar: common.VersionEnv,
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s,%s", common.Token, common.TokenShort),
			Value:  "",
			Usage:  "the OAuth token to use",
			EnvVar: common.TokenEnv,
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s,%s", common.ProjectID, common.ProjectIDShort),
			Value:  "",
			Usage:  "the project ID. your OAuth token must be authorized to access the project",
			EnvVar: common.ProjectIDEnv,
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s,%s", common.Host, common.HostShort),
			Value:  "",
			Usage:  "the host to use",
			EnvVar: common.HostEnv,
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "worker",
			Aliases:     []string{"w"},
			Usage:       "IronWorker commands",
			Subcommands: worker.Subcommands,
		},
		cli.Command{
			Name:        "queue",
			Aliases:     []string{"mq"},
			Usage:       "IronMQ commands",
			Subcommands: mq.Subcommands,
		},
	}
	app.Run(os.Args)
}
