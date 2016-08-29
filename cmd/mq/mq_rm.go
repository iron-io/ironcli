package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type MqRm struct {
	cli.Command
}

func NewMqRm(settings *config.Settings) *MqRm {
	mqRm := &MqRm{
		Command: cli.Command{
			Name:      "rm",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqRm
}

func (r MqRm) GetCmd() cli.Command {
	return r.Command
}
