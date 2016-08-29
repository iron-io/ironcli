package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type MqReverse struct {
	cli.Command
}

func NewMqReverse(settings *config.Settings) *MqReverse {
	mqReverse := &MqReverse{
		Command: cli.Command{
			Name:      "reverse",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqReverse
}

func (r MqReverse) GetCmd() cli.Command {
	return r.Command
}
