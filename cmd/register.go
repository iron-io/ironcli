package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/worker"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type Register struct {
	name            string
	config          string
	configFile      string
	maxConc         int
	retries         int
	retriesDelay    int
	defaultPriority int
	host            string
	codes           worker.Code

	cli.Command
}

func NewRegister(settings *config.Settings) *Register {
	register := &Register{}
	register.Command = cli.Command{
		Name:      "register",
		Usage:     "register worker in the project",
		UsageText: "doo - does the dooing",
		ArgsUsage: "[image] [command] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "name",
				Usage:       "name usage",
				Destination: &register.name,
			},
			cli.StringFlag{
				Name:        "config",
				Usage:       "config usage",
				Destination: &register.config,
			},
			cli.StringFlag{
				Name:        "configFile",
				Usage:       "configFile usage",
				Destination: &register.configFile,
			},
			cli.IntFlag{
				Name:        "maxConc",
				Usage:       "maxConc usage",
				Value:       -1,
				Destination: &register.maxConc,
			},
			cli.IntFlag{
				Name:        "retries",
				Usage:       "retries usage",
				Value:       0,
				Destination: &register.retries,
			},
			cli.IntFlag{
				Name:        "retriesDelay",
				Usage:       "retriesDelay usage",
				Value:       0,
				Destination: &register.retriesDelay,
			},
			cli.IntFlag{
				Name:        "defaultPriority",
				Usage:       "defaultPriority usage",
				Value:       -3,
				Destination: &register.defaultPriority,
			},
			cli.StringFlag{
				Name:        "host",
				Usage:       "host usage",
				Destination: &register.host,
			},
		},
		Action: func(c *cli.Context) error {
			err := register.Execute(c.Args().Tail(), c.Args().First())
			if err != nil {
				return err
			}

			if register.codes.Host != "" {
				fmt.Println(`Spinning up '` + register.codes.Name + `'`)
			} else {
				fmt.Println(`Registering worker '` + register.codes.Name + `'`)
			}

			code, err := common.PushCodes("", settings, register.codes)
			if err != nil {
				return err
			}

			if code.Host != "" {
				fmt.Println(`Hosted at: '` + code.Host + `'`)
			} else {
				fmt.Println(`Registered code package with id='` + code.Id + `'`)
			}

			return nil
		},
	}

	return register
}

func (r Register) GetCmd() cli.Command {
	return r.Command
}

func (r *Register) Execute(cmd []string, image string) error {
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
