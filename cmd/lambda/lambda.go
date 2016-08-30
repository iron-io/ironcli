package lambda

import (
	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type Lambda struct {
	cli.Command
}

func NewLambda(settings *config.Settings) *Lambda {
	lambda := &Lambda{
		Command: cli.Command{
			Name:      "lambda",
			Usage:     "manage lambda functions",
			ArgsUsage: "[command]",
			Subcommands: cli.Commands{
				NewLambdaCreateFunction().GetCmd(),
				NewLambdaAwsImport().GetCmd(),
				NewLambdaPublishFunction().GetCmd(),
				NewLambdaTestFunction().GetCmd(),
			},
		},
	}

	return lambda
}

func (r Lambda) GetCmd() cli.Command {
	return r.Command
}
