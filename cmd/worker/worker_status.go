package worker

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerStatus struct {
	wrkr common.Worker

	cli.Command
}

func NewWorkerStatus(settings *config.Settings) *WorkerStatus {
	workerStatus := &WorkerStatus{}

	workerStatus.Command = cli.Command{
		Name:      "status",
		Usage:     "get execution status of a task.",
		ArgsUsage: "[task_id]",
		Action: func(c *cli.Context) error {
			workerStatus.wrkr.Settings = *settings

			fmt.Println("LINES", `Getting status of task with id='`+c.Args().First()+`'`)

			taskInfo, err := workerStatus.wrkr.TaskInfo(c.Args().First())
			if err != nil {
				return err
			}

			fmt.Println("BLANKS", taskInfo.Status)

			return nil
		},
	}

	return workerStatus
}

func (r WorkerStatus) GetCmd() cli.Command {
	return r.Command
}
