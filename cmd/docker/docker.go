package docker

import "github.com/urfave/cli"

type Docker struct {
	cli.Command
}

func NewDocker() *Docker {
	docker := &Docker{
		Command: cli.Command{
			Name:      "docker",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Subcommands: cli.Commands{
				NewDockerLogin().GetCmd(),
			},
		},
	}

	return docker
}

func (r Docker) GetCmd() cli.Command {
	return r.Command
}
