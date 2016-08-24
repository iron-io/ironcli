package lambda

import "github.com/urfave/cli"

type Lambda struct {
	cli.Command
}

func NewLambda() *Lambda {
	lambda := &Lambda{
		Command: cli.Command{
			Name:      "lambda",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Subcommands: cli.Commands{
				NewLambdaCreateFunction().GetCmd(),
			},
		},
	}

	return lambda
}

func (r Lambda) GetCmd() cli.Command {
	return r.Command
}
