package worker

import (
	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type Worker struct {
	cli.Command
}

func NewWorker(settings *config.Settings) *Worker {
	worker := &Worker{
		Command: cli.Command{
			Name:      "worker",
			Usage:     "commands to interact with IronWorker.",
			ArgsUsage: "[command]",
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
