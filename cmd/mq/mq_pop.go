package mq

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type MqPop struct {
	cli.Command
}

func NewMqPop(settings *config.Settings) *MqPop {
	mqPop := &MqPop{
		Command: cli.Command{
			Name:      "pop",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqPop
}

func (r MqPop) GetCmd() cli.Command {
	return r.Command
}
