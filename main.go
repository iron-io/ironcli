package main

import (
	"fmt"
	"os"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/cmd"
	"github.com/iron-io/ironcli/cmd/docker"
	"github.com/iron-io/ironcli/cmd/lambda"
	"github.com/iron-io/ironcli/cmd/mq"
	"github.com/iron-io/ironcli/cmd/worker"
	"github.com/urfave/cli"
)

func main() {
	var (
		settings config.Settings
	)

	app := cli.NewApp()
	app.Name = "Iron CLI"
	app.Version = "0.3.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "iron.io",
			Email: "",
		},
	}
	app.HelpName = "iron"
	app.Usage = "Iron.io command line tools"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "project-id", Usage: "provide project ID"},
		cli.StringFlag{Name: "token", Usage: "provide OAuth token"},
		cli.StringFlag{Name: "env", Usage: "provide specific dev environment"},
	}

	// Init settings
	app.Before = func(c *cli.Context) error {
		// FIXME when someone will add a break changes for code
		if c.GlobalString("project-id") != "" {
			err := os.Setenv("IRON_PROJECT_ID", c.GlobalString("project-id"))
			if err != nil {
				return err
			}
		}

		// FIXME when someone will add a break changes for code
		if c.GlobalString("token") != "" {
			err := os.Setenv("IRON_TOKEN", c.GlobalString("token"))
			if err != nil {
				return err
			}
		}

		settings = config.ConfigWithEnv("iron_worker", c.GlobalString("env"))

		return nil
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		// FIXME when there will be a logger
		fmt.Fprintf(os.Stderr, "%q command not found.\n", command)
	}

	app.Commands = []cli.Command{
		cmd.NewRegister(&settings).GetCmd(),
		cmd.NewRun(&settings).GetCmd(),
		worker.NewWorker(&settings).GetCmd(),
		mq.NewMq(&settings).GetCmd(),
		docker.NewDocker(&settings).GetCmd(),
		lambda.NewLambda(&settings).GetCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		// FIXME when there will be a logger
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
