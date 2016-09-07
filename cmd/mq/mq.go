package mq

import (
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type Mq struct {
	cli.Command
}

func NewMq(settings *common.Settings) *Mq {
	mq := &Mq{
		Command: cli.Command{
			Name:      "mq",
			Usage:     "manage queues",
			ArgsUsage: "[command]",
			Before: func(c *cli.Context) error {
				settings.Product = "iron_mq"
				common.SetSettings(settings)

				return nil
			},
			Subcommands: cli.Commands{
				NewMqPush(settings).GetCmd(),
				NewMqClear(settings).GetCmd(),
				NewMqCreate(settings).GetCmd(),
				NewMqDelete(settings).GetCmd(),
				NewMqInfo(settings).GetCmd(),
				NewMqList(settings).GetCmd(),
				NewMqPeek(settings).GetCmd(),
				NewMqPop(settings).GetCmd(),
				NewMqReserve(settings).GetCmd(),
				NewMqRm(settings).GetCmd(),
			},
		},
	}

	return mq
}

func (r Mq) GetCmd() cli.Command {
	return r.Command
}
