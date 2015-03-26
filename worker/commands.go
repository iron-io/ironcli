package worker

import "github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"

var SubCommands = []cli.Command{
	[]cli.Command{
		Name:        "packages",
		Aliases:     []string{"pkg"},
		Usage:       "operate on IronWorker packages",
		SubCommands: packages.SubCommands,
	},
	[]cli.Command{
		Name:        "tasks",
		Aliases:     []string{"t"},
		Usage:       "operate on IronWorker tasks",
		SubCommands: tasks.SubCommands,
	},
	[]cli.Command{
		Name:       "scheduled",
		Aliases:    []string{"sched"},
		Usage:      "operate on IronWorker scheduled tasks",
		SubCommand: scheduled.SubCommands,
	},
	[]cli.Command{
		Name:       "stacks",
		Aliases:    []string{"stk"},
		Usage:      "operate on IronWorker stacks",
		SubCommand: stacks.SubCommands,
	},
}
