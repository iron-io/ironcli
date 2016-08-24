package worker

import (
	"fmt"

	"github.com/urfave/cli"
)

type WorkerSchedule struct {
	cli.Command
}

func NewWorkerSchedule() *WorkerSchedule {
	workerSchedule := &WorkerSchedule{
		Command: cli.Command{
			Name:      "schedule",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return workerSchedule
}

func (r WorkerSchedule) GetCmd() cli.Command {
	return r.Command
}
