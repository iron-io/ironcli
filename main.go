package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/iron-io/ironcli/cmd"
	"github.com/iron-io/ironcli/cmd/docker"
	"github.com/iron-io/ironcli/cmd/lambda"
	"github.com/iron-io/ironcli/cmd/mq"
	"github.com/iron-io/ironcli/cmd/worker"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

func main() {
	var (
		settings = &common.Settings{}
	)

	app := cli.NewApp()
	app.Name = "iron"
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
		// Setting old variables for color
		if runtime.GOOS == "windows" {
			common.Red = fmt.Sprint
			common.Yellow = fmt.Sprint
			common.Green = fmt.Sprint
		} else {
			common.Red = color.New(color.FgRed).SprintFunc()
			common.Yellow = color.New(color.FgYellow).SprintFunc()
			common.Green = color.New(color.FgGreen).SprintFunc()
		}

		settings.Env = c.GlobalString("env")

		if c.GlobalString("project-id") != "" {
			settings.Worker.ProjectId = c.GlobalString("project-id")
		}

		if c.GlobalString("token") != "" {
			settings.Worker.Token = c.GlobalString("token")
		}

		return nil
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		// FIXME when there will be a logger
		fmt.Fprintf(os.Stderr, "%q command not found.\n", command)
	}

	app.Commands = []cli.Command{
		cmd.NewRegister(settings).GetCmd(),
		cmd.NewRun(settings).GetCmd(),
		worker.NewWorker(settings).GetCmd(),
		mq.NewMq(settings).GetCmd(),
		docker.NewDocker(settings).GetCmd(),
		lambda.NewLambda(settings).GetCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		// FIXME when there will be a logger
		fmt.Fprint(os.Stderr, common.Red(common.BLANKS, err))
	}
}
