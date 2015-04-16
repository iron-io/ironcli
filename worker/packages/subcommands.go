package packages

import (
	"time"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/iron-io/iron_go3/worker"
	"github.com/iron-io/ironcli/common"
)

// cli.Command{
// 	Name:    "status",
// 	Aliases: []string{"stat"},
// 	Usage:   "get the status of an IronWorker task",
// 	Flags: []cli.Flag{
// 		cli.StringFlag{Name: "taskid", Usage: "the task ID to get the status of"},
// 	},
// 	Action: func(c *cli.Context) {},
// },
// cli.Command{
// 	Name:    "upload",
// 	Aliases: []string{"up"},
// 	Usage:   "upload a new IronWorker code",
// 	Flags:   uploadFlags,
// 	Action:  func(c *cli.Context) {},
// },

// var baseSchedUploadFlags = []cli.Flag{}
//
var uploadFlags = []cli.Flag{
	cli.StringFlag{Name: "payload", Usage: "worker payload"},
	cli.StringFlag{Name: "payloadfile", Usage: "worker payload file (payload takes precedence)"},
	cli.IntFlag{Name: "priority", Value: -1, Usage: "priority to run this code (0, 1 or 2)"},
	cli.IntFlag{Name: "timeout", Value: 3600, Usage: "max runtime for task in seconds (0 - 3600)"},
	cli.IntFlag{Name: "delay", Value: 0, Usage: "seconds to delay before queueing task"},
	cli.BoolFlag{Name: "wait", Usage: "wait for task to complete and print log"},
	cli.IntFlag{Name: "max-concurrency", Value: -1, Usage: "max workers to run in parallel"},
	cli.IntFlag{Name: "run-every", Value: -1, Usage: "time between runs, in seconds. -1 to run once"},
	cli.IntFlag{Name: "run-times", Value: 1, Usage: "number of times to run"},
	cli.StringFlag{Name: "end-at", Usage: "time or datetime of form 'Mon Jan 2 15:04:05 -0700 2006'"},
	cli.StringFlag{Name: "start-at", Usage: "time or datetime of form 'Mon Jan 2 15:04:05 -0700 2006'"},
	cli.IntFlag{Name: "retries", Value: 0, Usage: "max times to retry failed task. max 10"},
	cli.IntFlag{Name: "retries-delay", Value: 0, Usage: "time between retries, in seconds"},
	cli.StringFlag{Name: "config", Usage: "provide config string (JSON or YAML) that will be available in file on upload"},
	cli.StringFlag{Name: "stack", Value: "default", Usage: "the stack to run your codes in"},
}

var scheduleFlags = append(uploadFlags, []cli.Flag{
	cli.StringFlag{Name: "code-name", Usage: "the name of the code to schedule"},
	cli.StringFlag{Name: "name", Usage: "the name of the schedule"},
}...)

func getSchedule(g *common.GlobalFlags) worker.Schedule {
	delay := g.DurationOrFail("delay", time.Second, -1)
	endAt := g.TimeOrFail("end-at")
	maxConcurrency := g.IntOrFail("max-concurrency", -1)
	prio := g.IntOrFail("priority", -1)
	runEvery := g.IntOrFail("run-every", -2)
	runTimes := g.IntOrFail("run-times", -1)
	startAt := g.TimeOrFail("start-at")
	return worker.Schedule{
		CodeName:       g.StringOrFail("code-name"),
		Delay:          &delay,
		EndAt:          &endAt,
		MaxConcurrency: &maxConcurrency,
		Name:           g.StringOrFail("name"),
		Payload:        g.StringOrFail("payload"),
		Priority:       &prio,
		RunEvery:       &runEvery,
		RunTimes:       &runTimes,
		StartAt:        &startAt,
	}
}

var Subcommands = []cli.Command{
	cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List code packages",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "page",
				Value: -1,
				Usage: "the page number to get",
			},
			cli.IntFlag{
				Name:  "per-page",
				Usage: "the number of code packages to get per page",
				Value: 10,
			},
		},
		Action: common.WithGlobalFlags(list),
	},
	cli.Command{
		Name:    "upload",
		Aliases: []string{"u"},
		Usage:   "Upload or update a code package",
		Action:  common.WithGlobalFlags(upload),
	},
	cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "Get info about a code package",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "codeid", Usage: "the code ID to get info about"},
		},
		Action: common.WithGlobalFlags(info),
	},
	cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "Delete a code package",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "codeid", Usage: "the code ID to delete"},
		},
		Action: common.WithGlobalFlags(del),
	},
	cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "Download a code package",
		Action:  common.WithGlobalFlags(download),
	},
	cli.Command{
		Name:    "listrevs",
		Aliases: []string{"lsr"},
		Usage:   "List Code Package Revisions",
		Action:  common.WithGlobalFlags(listrevs),
	},
	cli.Command{
		Name:    "logs",
		Aliases: []string{"wl"},
		Usage:   "view IronWorker logs",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "taskid", Usage: "the task ID to get logs for"},
		},
		Action: common.WithGlobalFlags(logs),
	},
	cli.Command{
		Name:    "schedule",
		Aliases: []string{"sched"},
		Usage:   "schedule an IronWorker task to run",
		Flags:   scheduleFlags,
		Action:  common.WithGlobalFlags(schedule),
	},
}
