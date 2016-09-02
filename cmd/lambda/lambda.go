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
			Name: "lambda",
			Usage: `The Lambda commands allow packaging AWS Lambda compatible functions into Docker containers.
They also allow importing certain Lambda functions. Please see (link to either blog post or iron-io/lambda docs) for more information.`,
			ArgsUsage: "",
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
