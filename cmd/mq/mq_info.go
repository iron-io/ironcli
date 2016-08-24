package mq

import (
	"fmt"

	"github.com/urfave/cli"
)

type MqInfo struct {
	cli.Command
}

func NewMqInfo() *MqInfo {
	mqInfo := &MqInfo{
		Command: cli.Command{
			Name:      "info",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return mqInfo
}

func (r MqInfo) GetCmd() cli.Command {
	return r.Command
}
