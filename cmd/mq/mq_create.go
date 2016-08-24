package mq

import (
	"fmt"

	"github.com/urfave/cli"
)

type MqCreate struct {
	cli.Command
}

func NewMqCreate() *MqCreate {
	mqCreate := &MqCreate{
		Command: cli.Command{
			Name:      "creates",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqCreate
}

func (r MqCreate) GetCmd() cli.Command {
	return r.Command
}
