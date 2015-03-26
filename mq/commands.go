package mq

import "github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"

var SubCommands = []cli.Command{
	cli.Command{
		Name:    "enqueue",
		Aliases: []string{"enq"},
		Usage:   "enqueue a message",
		Action:  func(c *cli.Context) {},
	},
}
