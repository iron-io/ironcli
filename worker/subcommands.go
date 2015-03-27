package worker

import (
	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/worker/packages"
	"github.com/iron-io/ironcli/worker/scheduled"
	"github.com/iron-io/ironcli/worker/stacks"
	"github.com/iron-io/ironcli/worker/tasks"
)

var Subcommands = []cli.Command{
	cli.Command{
		Name:        "packages",
		Aliases:     []string{"pkg"},
		Usage:       "operate on IronWorker packages",
		Subcommands: packages.Subcommands,
	},
	cli.Command{
		Name:        "tasks",
		Aliases:     []string{"t"},
		Usage:       "operate on IronWorker tasks",
		Subcommands: tasks.Subcommands,
	},
	cli.Command{
		Name:        "scheduled",
		Aliases:     []string{"sched"},
		Usage:       "operate on IronWorker scheduled tasks",
		Subcommands: scheduled.Subcommands,
	},
	cli.Command{
		Name:        "stacks",
		Aliases:     []string{"stk"},
		Usage:       "operate on IronWorker stacks",
		Subcommands: stacks.Subcommands,
	},
}
