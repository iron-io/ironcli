package mq

import (
	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type Mq struct {
	cli.Command
}

func NewMq(settings *config.Settings) *Mq {
	mqSettings := config.ManualConfig("iron_mq", nil)
	mqSettings.Token = settings.Token
	mqSettings.ProjectId = settings.ProjectId

	mq := &Mq{
		Command: cli.Command{
			Name:      "mq",
			Usage:     "manage queues",
			ArgsUsage: "[command]",
			Subcommands: cli.Commands{
				NewMqPush(&mqSettings).GetCmd(),
				NewMqClear(&mqSettings).GetCmd(),
				NewMqCreate(&mqSettings).GetCmd(),
				NewMqDelete(&mqSettings).GetCmd(),
				NewMqInfo(&mqSettings).GetCmd(),
				NewMqList(&mqSettings).GetCmd(),
				NewMqPeek(&mqSettings).GetCmd(),
				NewMqPop(&mqSettings).GetCmd(),
				NewMqReserve(&mqSettings).GetCmd(),
				NewMqRm(&mqSettings).GetCmd(),
			},
		},
	}

	return mq
}

func (r Mq) GetCmd() cli.Command {
	return r.Command
}
