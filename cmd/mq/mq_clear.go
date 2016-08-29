package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type MqClear struct {
	cli.Command
}

func NewMqClear(settings *config.Settings) *MqClear {
	mqClear := &MqClear{
		Command: cli.Command{
			Name:      "clear",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqClear
}

func (r MqClear) GetCmd() cli.Command {
	return r.Command
}
