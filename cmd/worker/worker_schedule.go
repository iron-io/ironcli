package worker

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/iron-io/iron_go3/config"
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

func NewWorkerSchedule(settings *config.Settings) *WorkerSchedule {
	workerSchedule := &WorkerSchedule{}

	workerSchedule.Command = cli.Command{
		Name:  "schedule",
		Usage: "add worker as task on date to queue and run it",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "payload",
				Usage:       "",
				Destination: &workerSchedule.payload,
			},
			cli.StringFlag{
				Name:        "payload-file",
				Usage:       "",
				Destination: &workerSchedule.payloadFile,
			},
			cli.IntFlag{
				Name:        "priority",
				Usage:       "",
				Destination: &workerSchedule.priority,
			},
			cli.IntFlag{
				Name:        "timeout",
				Usage:       "",
				Destination: &workerSchedule.timeout,
			},
			cli.IntFlag{
				Name:        "delay",
				Usage:       "",
				Destination: &workerSchedule.delay,
			},
			cli.IntFlag{
				Name:        "run-every",
				Usage:       "",
				Destination: &workerSchedule.runEvery,
			},
			cli.IntFlag{
				Name:        "run-times",
				Usage:       "",
				Destination: &workerSchedule.runTimes,
			},
			cli.StringFlag{
				Name:        "cluster",
				Usage:       "",
				Destination: &workerSchedule.cluster,
			},
			cli.StringFlag{
				Name:        "label",
				Usage:       "",
				Destination: &workerSchedule.label,
			},
			cli.StringFlag{
				Name:        "start-at",
				Usage:       "",
				Destination: &workerSchedule.startAt,
			},
			cli.StringFlag{
				Name:        "end-at",
				Usage:       "",
				Destination: &workerSchedule.endAt,
			},
		},
		Action: func(c *cli.Context) error {
			workerSchedule.wrkr.Settings = *settings

			err := workerSchedule.Execute(c.Args().First())
			if err != nil {
				return err
			}

			fmt.Println("LINES", "Scheduling task '"+workerSchedule.sched.CodeName+"'")

			ids, err := workerSchedule.wrkr.Schedule(workerSchedule.sched)
			if err != nil {
				return err
			}

			id := ids[0]

			fmt.Printf("%s Scheduled task with id='%s'\n", "BLANKS", id)

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
