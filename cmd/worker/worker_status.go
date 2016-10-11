package worker

import (
	"fmt"

	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerStatus struct {
	Status string
	wrkr   common.Worker

	cli.Command
}

func NewWorkerStatus(settings *common.Settings) *WorkerStatus {
	workerStatus := &WorkerStatus{}

	workerStatus.Command = cli.Command{
		Name:      "status",
		Usage:     "get execution status of a task.",
		ArgsUsage: "[task_id]",
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := workerStatus.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return workerStatus
}

func (w WorkerStatus) GetCmd() cli.Command {
	return w.Command
}

func (w *WorkerStatus) Action(taskID string, settings *common.Settings) error {
	w.wrkr.Settings = settings.Worker

	fmt.Println(common.LINES, `Getting status of task with id='`+taskID+`'`)

	taskInfo, err := w.wrkr.TaskInfo(taskID)
	if err != nil {
		return err
	}

	fmt.Println(common.BLANKS, taskInfo.Status)

	w.Status = taskInfo.Status

	return nil
}
