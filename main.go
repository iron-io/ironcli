package main

import (
	"fmt"
	"os"
	"time"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/cmd"
	"github.com/iron-io/ironcli/cmd/docker"
	"github.com/iron-io/ironcli/cmd/lambda"
	"github.com/iron-io/ironcli/cmd/mq"
	"github.com/iron-io/ironcli/cmd/worker"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

type genericType struct {
	s string
	*config.Settings
}

func (g *genericType) Set(value string) error {
	g.s = value
	return nil
}

func (g *genericType) String() string {
	return g.s
}

func main() {
	app := cli.NewApp()
	app.Name = "Iron CLI"
	app.Version = "0.3.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "iron.io",
			Email: "",
		},
	}
	app.HelpName = "iron"
	app.Usage = "Go version of the Iron.io command line tools"

	app.Flags = []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.token"}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.project_id"}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.host"}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.scheme"}),
		altsrc.NewIntFlag(cli.IntFlag{Name: "settings.port"}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.api_version"}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "settings.user_agent"}),
		cli.StringFlag{Name: "project-id"},
		cli.StringFlag{Name: "token"},
		cli.StringFlag{Name: "config", Value: "env/config.yaml"},
	}

	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))

	// Init settings
	app.After = func(c *cli.Context) error {
		settings := config.Settings{
			Token:      c.GlobalString("settings.token"),
			ProjectId:  c.GlobalString("settings.project_id"),
			Host:       c.GlobalString("settings.host"),
			Scheme:     c.GlobalString("settings.scheme"),
			Port:       uint16(c.GlobalInt("settings.port")),
			ApiVersion: c.GlobalString("settings.api_version"),
			UserAgent:  c.GlobalString("settings.user_agent"),
		}

		if c.GlobalString("project-id") != "" {
			settings.ProjectId = c.GlobalString("project-id")
		}

		if c.GlobalString("token") != "" {
			settings.Token = c.GlobalString("token")
		}

		return nil
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("%q command not found.\n", command)
	}

	app.Commands = []cli.Command{
		cmd.NewRegister().GetCmd(),
		cmd.NewRun().GetCmd(),
		worker.NewWorker().GetCmd(),
		mq.NewMq().GetCmd(),
		docker.NewDocker().GetCmd(),
		lambda.NewLambda().GetCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("WRONG: %#v\n", err)
	}
}
