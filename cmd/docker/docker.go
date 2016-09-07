package docker

import (
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type Docker struct {
	cli.Command
}

func NewDocker(settings *common.Settings) *Docker {
	docker := &Docker{
		Command: cli.Command{
			Name:  "docker",
			Usage: "manage Docker credentials.",
			Before: func(c *cli.Context) error {
				settings.Product = "iron_worker"
				common.SetSettings(settings)

				return nil
			},
			Subcommands: cli.Commands{
				NewDockerLogin(settings).GetCmd(),
			},
		},
	}

	return docker
}

func (r Docker) GetCmd() cli.Command {
	return r.Command
}
