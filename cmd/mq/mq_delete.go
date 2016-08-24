package mq

import (
	"fmt"

	"github.com/urfave/cli"
)

type MqDelete struct {
	cli.Command
}

func NewMqDelete() *MqDelete {
	mqDelete := &MqDelete{
		Command: cli.Command{
			Name:      "delete",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqDelete
}

func (r MqDelete) GetCmd() cli.Command {
	return r.Command
}
