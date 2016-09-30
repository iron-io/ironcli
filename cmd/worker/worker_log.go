package worker

import (
	"fmt"

	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerLog struct {
	Log  string
	wrkr common.Worker

	cli.Command
}

func NewWorkerLog(settings *common.Settings) *WorkerLog {
	workerLog := &WorkerLog{}

	workerLog.Command = cli.Command{
		Name:      "log",
		Usage:     "get log output of a task that has finished executing.",
		ArgsUsage: "[task-id]",
		Action: func(c *cli.Context) error {
			err := workerLog.Action(c.Args().First(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return workerLog
}

func (w WorkerLog) GetCmd() cli.Command {
	return w.Command
}

func (w *WorkerLog) Action(taskID string, settings *common.Settings) error {
	w.wrkr.Settings = settings.Worker

	fmt.Println(common.LINES, "Getting log for task with id='"+taskID+"'")

	out, err := w.wrkr.TaskLog(taskID)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	w.Log = string(out)

	return nil
}
