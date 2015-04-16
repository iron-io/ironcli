package mq

import (
	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/common"
)

var qFlag = cli.StringFlag{
	Name:  "queue",
	Usage: "the name of the queue to operate on",
}

var Subcommands = []cli.Command{
	cli.Command{
		Name:    "push",
		Aliases: []string{"p"},
		Usage:   "push a message",
		Flags: []cli.Flag{
			qFlag,
			cli.IntFlag{
				Name:  "delay",
				Value: 0,
				Usage: "time in seconds to wait before enqueue",
			},
			cli.StringFlag{
				Name:  "body",
				Usage: "body of the message",
			},
		},
		Action: common.WithGlobalFlags(enqueue),
	},
	cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete a message",
		Flags: []cli.Flag{
			qFlag,
			cli.StringFlag{
				Name:  "id",
				Usage: "the message ID",
			},
			cli.StringFlag{
				Name:  "reservation_id",
				Usage: "the message's reservation ID",
			},
		},
		Action: common.WithGlobalFlags(delete),
	},
}
