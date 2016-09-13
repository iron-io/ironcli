package worker

import (
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type Worker struct {
	cli.Command
}

func NewWorker(settings *common.Settings) *Worker {
	worker := &Worker{
		Command: cli.Command{
			Name:      "worker",
			Usage:     "commands to interact with IronWorker.",
			ArgsUsage: "[command]",
			Before: func(c *cli.Context) error {
				settings.Product = "iron_worker"
				if err := common.SetSettings(settings); err != nil {
					return err
				}

				return nil
			},
			Subcommands: cli.Commands{
				NewWorkerUpload(settings).GetCmd(),
				NewWorkerLog(settings).GetCmd(),
				NewWorkerQueue(settings).GetCmd(),
				NewWorkerSchedule(settings).GetCmd(),
				NewWorkerStatus(settings).GetCmd(),
			},
		},
	}

	return worker
}

func (r Worker) GetCmd() cli.Command {
	return r.Command
}
