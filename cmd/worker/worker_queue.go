package worker

import (
	"fmt"

	"github.com/urfave/cli"
)

type WorkerQueue struct {
	cli.Command
}

func NewWorkerQueue() *WorkerQueue {
	workerQueue := &WorkerQueue{
		Command: cli.Command{
			Name:      "queue",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return workerQueue
}

func (r WorkerQueue) GetCmd() cli.Command {
	return r.Command
}
