package worker

import (
	"fmt"

	"github.com/urfave/cli"
)

type WorkerUpload struct {
	cli.Command
}

func NewWorkerUpload() *WorkerUpload {
	workerUpload := &WorkerUpload{
		Command: cli.Command{
			Name:      "upload",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return workerUpload
}

func (r WorkerUpload) GetCmd() cli.Command {
	return r.Command
}
