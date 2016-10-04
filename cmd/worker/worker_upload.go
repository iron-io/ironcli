package worker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/iron-io/iron_go3/worker"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type WorkerUpload struct {
	name            string
	config          string
	configFile      string
	maxConc         int
	retries         int
	retriesDelay    int
	defaultPriority int
	Zip             string
	host            string
	codes           worker.Code
	WorkerID        string

	cli.Command
}

func NewWorkerUpload(settings *common.Settings) *WorkerUpload {
	workerUpload := &WorkerUpload{}
	workerUpload.Command = cli.Command{
		Name:      "upload",
		Usage:     "upload a new code package to IronWorker.",
		ArgsUsage: "[image]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "name",
				Usage:       "override code package name",
				Destination: &workerUpload.name,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "provide config string (re: JSON/YAML) that will be available in file on upload",
				Destination: &workerUpload.config,
			},
			cli.StringFlag{
				Name:        "config-file",
				Usage:       "upload file for worker config",
				Destination: &workerUpload.configFile,
			},
			cli.IntFlag{
				Name:        "max-conc",
				Usage:       "max workers to run in parallel. default is no limit",
				Value:       -1,
				Destination: &workerUpload.maxConc,
			},
			cli.IntFlag{
				Name:        "retries",
				Usage:       "max times to retry failed task, max 10, default 0",
				Value:       0,
				Destination: &workerUpload.retries,
			},
			cli.IntFlag{
				Name:        "retries-delay",
				Usage:       "time between retries, in seconds. default 0",
				Value:       0,
				Destination: &workerUpload.retriesDelay,
			},
			cli.IntFlag{
				Name:        "default-priority",
				Usage:       "0(default), 1 or 2",
				Value:       -3,
				Destination: &workerUpload.defaultPriority,
			},
			cli.StringFlag{
				Name:        "zip",
				Usage:       "optional: name of zip file where code resides",
				Destination: &workerUpload.Zip,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "paas host",
				Destination: &workerUpload.host,
			},
		},
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := workerUpload.Action(c.Args().First(), c.Args().Tail(), settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return workerUpload
}

func (w WorkerUpload) GetCmd() cli.Command {
	return w.Command
}

func (w *WorkerUpload) Execute(cmd []string, image string) error {
	w.codes.Command = strings.TrimSpace(strings.Join(cmd, " "))
	w.codes.Image = image
	w.codes.Name = image

	if w.name != "" {
		w.codes.Name = w.name
	} else {
		w.codes.Name = w.codes.Image
		if strings.ContainsRune(w.codes.Name, ':') {
			arr := strings.SplitN(w.codes.Name, ":", 2)
			w.codes.Name = arr[0]
		}
	}

	if w.Zip != "" {
		if !strings.HasSuffix(w.Zip, ".zip") {
			return errors.New("file extension must be .zip, got: " + w.Zip)
		}

		if _, err := os.Stat(w.Zip); err != nil {
			return err
		}
	}

	w.codes.MaxConcurrency = w.maxConc
	w.codes.Retries = w.retries
	w.codes.RetriesDelay = w.retriesDelay
	w.codes.Config = w.config
	w.codes.DefaultPriority = w.defaultPriority

	if w.host != "" {
		w.codes.Host = w.host
	}

	if w.configFile != "" {
		pload, err := ioutil.ReadFile(w.configFile)
		if err != nil {
			return err
		}
		w.codes.Config = string(pload)
	}

	return nil
}

func (w *WorkerUpload) Action(image string, cmd []string, settings *common.Settings) error {
	err := w.Execute(cmd, image)
	if err != nil {
		return err
	}

	if w.codes.Host != "" {
		fmt.Println(common.LINES, `Spinning up '`+w.codes.Name+`'`)
	} else {
		fmt.Println(common.LINES, `Uploading worker '`+w.codes.Name+`'`)
	}

	code, err := common.PushCodes(w.Zip, &settings.Worker, w.codes)
	if err != nil {
		return err
	}

	fmt.Println(code)

	if code.Host != "" {
		fmt.Println(common.BLANKS, common.Green(`Hosted at: '`+code.Host+`'`))
	} else {
		fmt.Println(common.BLANKS, common.Green(`Uploaded code package with id='`+code.Id+`'`))
	}

	fmt.Println(common.BLANKS, common.Green(settings.HUDUrlStr+"codes/"+code.Id+common.INFO))

	w.WorkerID = code.Id

	return nil
}
