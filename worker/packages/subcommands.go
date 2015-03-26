package packages

import "github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"

// cli.Command{
// 	Name:    "logs",
// 	Aliases: []string{"wl"},
// 	Usage:   "view IronWorker logs",
// 	Flags:   []cli.Flag{Name: "taskid", Usage: "the task ID to get logs for"},
// 	Action:  func(c *cli.Context) {},
// },
// cli.Command{
// 	Name:    "schedule",
// 	Aliases: []string{"sched"},
// 	Usage:   "schedule an IronWorker task to run",
// 	Flags:   uploadFlags,
// 	Action:  func(c *cli.Context) {},
// },
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
// var uploadFlags = []cli.Flag{
// 	cli.StringFlag{Name: "payload", Usage: "worker payload"},
// 	cli.StringFlag{Name: "payloadfile", Usage: "worker payload file (payload takes precedence)"},
// 	cli.IntFlag{Name: "priority", Value: 0, Usage: "priority to run this code (0, 1 or 2)"},
// 	cli.IntFlag{Name: "timeout", Value: 3600, Usage: "max runtime for task in seconds (0 - 3600)"},
// 	cli.IntFlag{Name: "delay", Value: 0, Usage: "seconds to delay before queueing task"},
// 	cli.BoolFlag{Name: "wait", Value: false, Usage: "wait for task to complete and print log"},
// 	cli.IntFlag{Name: "max-concurrency", Value: -1, Usage: "max workers to run in parallel"},
// 	cli.IntFlag{Name: "run-every", Value: -1, Usage: "time between runs, in seconds. -1 to run once"},
// 	cli.IntFlag{Name: "run-times", Value: 1, Usage: "number of times to run"},
// 	cli.StringFlag{Name: "end-at", Usage: "time or datetime of form 'Mon Jan 2 15:04:05 -0700 2006'"},
// 	cli.StringFlag{Name: "start-at", Usage: "time or datetime of form 'Mon Jan 2 15:04:05 -0700 2006'"},
// 	cli.StringFlag{Name: "retries", Value: 0, Usage: "max times to retry failed task. max 10"},
// 	cli.IntFlag{Name: "retries-delay", Value: 0, Usage: "time between retries, in seconds"},
// 	cli.StringFlag{Name: "config", Usage: "provide config string (JSON or YAML) that will be available in file on upload"},
// 	cli.StringFlag{Name: "stack", Value: "default", Usage: "the stack to run your codes in"},
// }

var SubCommands = []cli.Command{
	[]cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List code packages",
		Action:  list,
	},
	[]cli.Command{
		Name:    "upload",
		Aliases: []string{"u"},
		Usage:   "Upload or update a code package",
		Action:  upload,
	},
	[]cli.Command{
		Name:    "info",
		Aliases: []string{"i"},
		Usage:   "Get info about a code package",
		Action:  info,
	},
	[]cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "Delete a code package",
		Action:  del,
	},
	[]cli.Command{
		Name:    "download",
		Aliases: []string{"dl"},
		Usage:   "Download a code package",
		Action:  download,
	},
	[]cli.Command{
		Name:    "listrevs",
		Aliases: []string{"lsr"},
		Usage:   "List Code Package Revisions",
		Action:  listrevs,
	},
	[]cli.Command{
		Name:    "pause",
		Aliases: []string{"p"},
		Usage:   "Pause Task Queue for Code Package",
		Action:  pause,
	},
	[]cli.Command{
		Name:    "resume",
		Aliases: []string{"r"},
		Usage:   "Resume Paused Task Queue for Code Package",
		Action:  resume,
	},
}
