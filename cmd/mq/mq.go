package mq

import "github.com/urfave/cli"

type Mq struct {
	cli.Command
}

func NewMq() *Mq {
	mq := &Mq{
		Command: cli.Command{
			Name:      "mq",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Subcommands: cli.Commands{
				NewMqPush().GetCmd(),
			},
		},
	}

	return mq
}

func (r Mq) GetCmd() cli.Command {
	return r.Command
}
