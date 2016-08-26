package worker

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerLog struct {
	wrkr common.Worker

	cli.Command
}

func NewWorkerLog(settings *config.Settings) *WorkerLog {
	workerLog := &WorkerLog{}

	workerLog.Command = cli.Command{
		Name:      "log",
		Usage:     "get log from worker of queue",
		ArgsUsage: "[task-id]",
		Action: func(c *cli.Context) error {
			workerLog.wrkr.Settings = *settings

			fmt.Println("LINES", "Getting log for task with id='"+c.Args().First()+"'")
			out, err := workerLog.wrkr.TaskLog(c.Args().First())
			if err != nil {
				return err
			}

			fmt.Println(string(out))

			return nil
		},
	}

	return workerLog
}

func (r WorkerLog) GetCmd() cli.Command {
	return r.Command
}
