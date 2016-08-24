package worker

import (
	"fmt"

	"github.com/urfave/cli"
)

type WorkerLog struct {
	cli.Command
}

func NewWorkerLog() *WorkerLog {
	workerLog := &WorkerLog{
		Command: cli.Command{
			Name:      "log",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return workerLog
}

func (r WorkerLog) GetCmd() cli.Command {
	return r.Command
}
