package mq

import (
	"fmt"

	"github.com/urfave/cli"
)

type MqList struct {
	cli.Command
}

func NewMqList() *MqList {
	mqList := &MqList{
		Command: cli.Command{
			Name:      "list",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqList
}

func (r MqList) GetCmd() cli.Command {
	return r.Command
}
