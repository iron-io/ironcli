package worker

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerSchedule struct {
	payload     string
	payloadFile string
	priority    int
	timeout     int
	delay       int
	maxConc     int
	runEvery    int
	runTimes    int
	cluster     string
	endAt       string
	startAt     string
	label       string
	sched       common.Schedule
	wrkr        common.Worker

	cli.Command
}

func NewWorkerSchedule(settings *common.Settings) *WorkerSchedule {
	workerSchedule := &WorkerSchedule{}

	workerSchedule.Command = cli.Command{
		Name:  "schedule",
		Usage: "schedule a new task to run at a specified time.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "payload",
				Usage:       "give worker payload",
				Destination: &workerSchedule.payload,
			},
			cli.StringFlag{
				Name:        "payload-file",
				Usage:       "give worker payload of file contents",
				Destination: &workerSchedule.payloadFile,
			},
			cli.IntFlag{
				Name:        "priority",
				Usage:       "0(default), 1 or 2; uses worker's default priority if unset",
				Value:       -3,
				Destination: &workerSchedule.priority,
			},
			cli.IntFlag{
				Name:        "timeout",
				Usage:       "0-3600(default) max runtime for task in seconds",
				Value:       3600,
				Destination: &workerSchedule.timeout,
			},
			cli.IntFlag{
				Name:        "delay",
				Usage:       "seconds to delay before queueing task",
				Destination: &workerSchedule.delay,
			},
			cli.IntFlag{
				Name:        "run-every",
				Usage:       "time between runs in seconds (>= 60), default is run once",
				Value:       -1,
				Destination: &workerSchedule.runEvery,
			},
			cli.IntFlag{
				Name:        "run-times",
				Usage:       "number of times a task will run",
				Value:       0,
				Destination: &workerSchedule.runTimes,
			},
			cli.StringFlag{
				Name:        "cluster",
				Usage:       "optional: specify cluster to queue task on",
				Destination: &workerSchedule.cluster,
			},
			cli.StringFlag{
				Name:        "label",
				Usage:       "optional: specify label for a task",
				Destination: &workerSchedule.label,
			},
			cli.StringFlag{
				Name:        "start-at",
				Usage:       "time or datetime in RFC3339 format: '2006-01-02T15:04:05Z07:00'",
				Destination: &workerSchedule.startAt,
			},
			cli.StringFlag{
				Name:        "end-at",
				Usage:       "time or datetime in RFC3339 format: '2006-01-02T15:04:05Z07:00'",
				Destination: &workerSchedule.endAt,
			},
		},
		Action: func(c *cli.Context) error {
			workerSchedule.wrkr.Settings = settings.Worker

			err := workerSchedule.Execute(c.Args().First())
			if err != nil {
				return err
			}

			fmt.Println(common.LINES, "Scheduling task '"+workerSchedule.sched.CodeName+"'")

			ids, err := workerSchedule.wrkr.Schedule(workerSchedule.sched)
			if err != nil {
				return err
			}

			id := ids[0]

			fmt.Printf("%s Scheduled task with id='%s'\n", common.BLANKS, id)
			fmt.Println(common.BLANKS, settings.HUDUrlStr+"scheduled_tasks/"+id+common.INFO)

			return nil
		},
	}

	return workerSchedule
}

func (r WorkerSchedule) GetCmd() cli.Command {
	return r.Command
}

func (r *WorkerSchedule) Execute(codePackageName string) error {
	delay := time.Duration(r.delay) * time.Second

	var priority *int
	if r.priority > -3 && r.priority < 3 {
		priority = &r.priority
	}

	r.sched = common.Schedule{
		CodeName: codePackageName,
		Delay:    &delay,
		Priority: priority,
		RunTimes: &r.runTimes,
		Cluster:  r.cluster,
		Label:    r.label,
	}

	payload := r.payload
	if r.payloadFile != "" {
		pload, err := ioutil.ReadFile(r.payloadFile)
		if err != nil {
			return err
		}

		payload = string(pload)
	}

	if payload != "" {
		r.sched.Payload = payload
	} else {
		r.sched.Payload = "{}" // if we don't set this, it gets a 400 from API.
	}

	if r.endAt != "" {
		t, _ := time.Parse(time.RFC3339, r.endAt) // checked in validateFlags()
		r.sched.EndAt = &t
	}
	if r.startAt != "" {
		t, _ := time.Parse(time.RFC3339, r.startAt)
		r.sched.StartAt = &t
	}
	if r.maxConc > 0 {
		r.sched.MaxConcurrency = &r.maxConc
	}
	if r.runEvery > 0 {
		r.sched.RunEvery = &r.runEvery
	}

	return nil
}
