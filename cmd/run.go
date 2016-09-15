package cmd

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

type Run struct {
	name            string
	config          string
	configFile      string
	maxConc         int
	retries         int
	retriesDelay    int
	defaultPriority int
	zip             string
	host            string
	codes           worker.Code

	cli.Command
}

func NewRun(settings *common.Settings) *Run {
	run := &Run{}

	run.Command = cli.Command{
		Name:      "run",
		Usage:     "run a new task.",
		ArgsUsage: "[image] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "name",
				Usage:       "override code package name",
				Destination: &run.name,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "provide config string (re: JSON/YAML) that will be available in file on upload",
				Destination: &run.config,
			},
			cli.StringFlag{
				Name:        "config-file",
				Usage:       "upload file for worker config",
				Destination: &run.configFile,
			},
			cli.IntFlag{
				Name:        "max-conc",
				Usage:       "max workers to run in parallel. default is no limit",
				Value:       -1,
				Destination: &run.maxConc,
			},
			cli.IntFlag{
				Name:        "retries",
				Usage:       "max times to retry failed task, max 10, default 0",
				Value:       0,
				Destination: &run.retries,
			},
			cli.IntFlag{
				Name:        "retries-delay",
				Usage:       "time between retries, in seconds. default 0",
				Value:       0,
				Destination: &run.retriesDelay,
			},
			cli.IntFlag{
				Name:        "default-priority",
				Usage:       "0(default), 1 or 2",
				Value:       -3,
				Destination: &run.defaultPriority,
			},
			cli.StringFlag{
				Name:        "zip",
				Usage:       "optional: name of zip file where code resides",
				Destination: &run.zip,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "paas host",
				Destination: &run.host,
			},
		},
		Before: func(c *cli.Context) error {
			settings.Product = "iron_worker"
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := run.Execute(c.Args().Tail(), c.Args().First())
			if err != nil {
				return err
			}

			if run.codes.Host != "" {
				fmt.Println(common.LINES, `Spinning up '`+run.codes.Name+`'`)
			} else {
				fmt.Println(common.LINES, `Uploading worker '`+run.codes.Name+`'`)
			}

			code, err := common.PushCodes(run.zip, &settings.Worker, run.codes)
			if err != nil {
				return err
			}

			if code.Host != "" {
				fmt.Println(common.BLANKS, common.Green(`Hosted at: '`+code.Host+`'`))
			} else {
				fmt.Println(common.BLANKS, common.Green(`Uploaded code package with id='`+code.Id+`'`))
			}

			fmt.Println(common.BLANKS, common.Green(settings.HUDUrlStr+"codes/"+code.Id+common.INFO))

			return nil
		},
	}

	return run
}

func (r Run) GetCmd() cli.Command {
	return r.Command
}

func (r *Run) Execute(cmd []string, image string) error {
	r.codes.Command = strings.TrimSpace(strings.Join(cmd, " "))
	r.codes.Image = image
	r.codes.Name = image

	if r.name != "" {
		r.codes.Name = r.name
	} else {
		r.codes.Name = r.codes.Image
		if strings.ContainsRune(r.codes.Name, ':') {
			arr := strings.SplitN(r.codes.Name, ":", 2)
			r.codes.Name = arr[0]
		}
	}

	if r.zip != "" {
		if !strings.HasSuffix(r.zip, ".zip") {
			return errors.New("file extension must be .zip, got: " + r.zip)
		}

		if _, err := os.Stat(r.zip); err != nil {
			return err
		}
	}

	r.codes.MaxConcurrency = r.maxConc
	r.codes.Retries = r.retries
	r.codes.RetriesDelay = r.retriesDelay
	r.codes.Config = r.config
	r.codes.DefaultPriority = r.defaultPriority

	if r.host != "" {
		r.codes.Host = r.host
	}

	if r.configFile != "" {
		pload, err := ioutil.ReadFile(r.configFile)
		if err != nil {
			return err
		}
		r.codes.Config = string(pload)
	}

	return nil
}
