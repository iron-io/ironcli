package worker

import (
	"fmt"

	"github.com/urfave/cli"
)

type WorkerStatus struct {
	cli.Command
}

func NewWorkerStatus() *WorkerStatus {
	workerStatus := &WorkerStatus{
		Command: cli.Command{
			Name:      "status",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return workerStatus
}

func (r WorkerStatus) GetCmd() cli.Command {
	return r.Command
}
