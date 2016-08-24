package worker

import "github.com/urfave/cli"

type Worker struct {
	cli.Command
}

func NewWorker() *Worker {
	worker := &Worker{
		Command: cli.Command{
			Name:      "worker",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Subcommands: cli.Commands{
				NewWorkerUpload().GetCmd(),
				NewWorkerLog().GetCmd(),
				NewWorkerQueue().GetCmd(),
				NewWorkerSchedule().GetCmd(),
				NewWorkerStatus().GetCmd(),
			},
		},
	}

	return worker
}

func (r Worker) GetCmd() cli.Command {
	return r.Command
}
