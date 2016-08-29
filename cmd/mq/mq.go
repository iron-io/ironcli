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
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Subcommands: cli.Commands{
				NewMqPush(&mqSettings).GetCmd(),
				NewMqClear().GetCmd(),
				NewMqCreate(&mqSettings).GetCmd(),
				NewMqDelete().GetCmd(),
				NewMqInfo().GetCmd(),
				NewMqList().GetCmd(),
				NewMqPeek().GetCmd(),
				NewMqPop().GetCmd(),
				NewMqReverse().GetCmd(),
				NewMqRm().GetCmd(),
			},
		},
	}

	return mq
}

func (r Mq) GetCmd() cli.Command {
	return r.Command
}
