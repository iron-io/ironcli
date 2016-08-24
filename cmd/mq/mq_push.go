package mq

import (
	"fmt"

	"github.com/urfave/cli"
)

type MqPush struct {
	cli.Command
}

func NewMqPush() *MqPush {
	mqPush := &MqPush{
		Command: cli.Command{
			Name:      "push",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqPush
}

func (r MqPush) GetCmd() cli.Command {
	return r.Command
}
