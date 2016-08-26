package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerQueue struct {
	payload           string
	payloadFile       string
	priority          int
	timeout           int
	delay             int
	wait              bool
	cluster           string
	label             string
	encryptionKey     string
	encryptionKeyFile string
	task              common.Task
	wrkr              common.Worker

	cli.Command
}

func NewWorkerQueue(settings *config.Settings) *WorkerQueue {
	workerQueue := &WorkerQueue{}

	workerQueue.Command = cli.Command{
		Name:      "queue",
		Usage:     "do the doo",
		UsageText: "doo - does the dooing",
		ArgsUsage: "[image] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "payload",
				Usage:       "payload usage",
				Destination: &workerQueue.payload,
			},
			cli.StringFlag{
				Name:        "payload-file",
				Usage:       "payload-file usage",
				Destination: &workerQueue.payloadFile,
			},
			cli.IntFlag{
				Name:        "priority",
				Usage:       "priority usage",
				Destination: &workerQueue.priority,
			},
			cli.IntFlag{
				Name:        "timeout",
				Usage:       "timeout usage",
				Destination: &workerQueue.timeout,
			},
			cli.IntFlag{
				Name:        "delay",
				Usage:       "delay usage",
				Destination: &workerQueue.delay,
			},
			cli.BoolFlag{
				Name:        "wait",
				Usage:       "wait usage",
				Destination: &workerQueue.wait,
			},
			cli.StringFlag{
				Name:        "cluster",
				Usage:       "cluster usage",
				Destination: &workerQueue.cluster,
			},
			cli.StringFlag{
				Name:        "label",
				Usage:       "label usage",
				Destination: &workerQueue.label,
			},
			cli.StringFlag{
				Name:        "encryption-key",
				Usage:       "encryption-key usage",
				Destination: &workerQueue.encryptionKey,
			},
			cli.StringFlag{
				Name:        "encryption-key-file",
				Usage:       "encryption-key-file usage",
				Destination: &workerQueue.encryptionKeyFile,
			},
		},
		Action: func(c *cli.Context) error {
			workerQueue.wrkr.Settings = *settings

			err := workerQueue.Execute(c.Args().First())
			if err != nil {
				return err
			}

			ids, err := workerQueue.wrkr.TaskQueue(workerQueue.task)
			if err != nil {
				return err
			}
			id := ids[0]

			fmt.Printf("%s Queued task with id='%s'\n", "BLANKS", id)

			if workerQueue.wait {
				fmt.Println("LINES", "Waiting for task to start running")

				done := make(chan struct{})
				go workerQueue.runWatch(done, "queued")
				workerQueue.waitForRunning(id)
				close(done)

				// TODO print actual queued time?
				fmt.Println("LINES", "Task running, waiting for completion")

				done = make(chan struct{})
				go workerQueue.runWatch(done, "running")
				ti := <-workerQueue.wrkr.WaitForTask(id)
				close(done)
				if ti.Msg != "" {
					return fmt.Errorf("error running task: %v", ti.Msg)
				}

				log, err := workerQueue.wrkr.TaskLog(id)
				if err != nil {
					return fmt.Errorf("error getting log: %v", err)
				}

				fmt.Println("LINES", "Done")
				fmt.Println("LINES", "Printing Log:")
				fmt.Printf("%s", string(log))
			}

			return nil
		},
	}

	return workerQueue
}

func (r WorkerQueue) GetCmd() cli.Command {
	return r.Command
}

func (r *WorkerQueue) waitForRunning(taskId string) {
	for {
		info, err := r.wrkr.TaskInfo(taskId)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error getting task info:", err)
			return
		}

		if info.Status == "queued" {
			time.Sleep(100 * time.Millisecond)
		} else {
			return
		}
	}
}

func (r *WorkerQueue) runWatch(done <-chan struct{}, state string) {
	start := time.Now()
	var elapsed time.Duration
	var h, m, s, ms int64
	for {
		select {
		case <-time.After(time.Millisecond):
		case <-done:
			fmt.Println("LINES", state+":", fmt.Sprintf("%v:%v:%v:%v\r", h, m, s, ms))
			return
		}
		elapsed = time.Since(start)

		h = common.Mod(elapsed.Hours(), 24)
		m = common.Mod(elapsed.Minutes(), 60)
		s = common.Mod(elapsed.Seconds(), 60)
		ms = common.Mod(float64(elapsed.Nanoseconds())/1000, 100)

		fmt.Println("LINES", " "+state+":", fmt.Sprintf(" %v:%v:%v:%v\r", h, m, s, ms))
	}
}

func (r *WorkerQueue) Execute(codePackageName string) error {
	payload := r.payload
	if r.payloadFile != "" {
		pload, err := ioutil.ReadFile(r.payloadFile)
		if err != nil {
			return err
		}
		payload = string(pload)
	}

	delay := time.Duration(r.delay) * time.Second
	timeout := time.Duration(r.timeout) * time.Second

	var priority int = -3
	if r.priority > -3 && r.priority < 3 {
		priority = r.priority
	}

	encryptionKey := []byte(r.encryptionKey)
	if r.encryptionKeyFile != "" {
		var err error
		encryptionKey, err = ioutil.ReadFile(r.encryptionKeyFile)
		if err != nil {
			return err
		}
	}

	if len(encryptionKey) > 0 {
		var err error
		payload, err = common.RsaEncrypt(encryptionKey, payload)
		if err != nil {
			return err
		}
	}

	r.task = common.Task{
		CodeName: codePackageName,
		Payload:  payload,
		Priority: priority,
		Timeout:  &timeout,
		Delay:    &delay,
		Cluster:  r.cluster,
		Label:    r.label,
	}

	return nil
}
