package lambda

import (
	"fmt"

	"github.com/urfave/cli"
)

type LambdaAwsImport struct {
	cli.Command
}

func NewLambdaAwsImport() *LambdaAwsImport {
	lambdaAwsImport := &LambdaAwsImport{
		Command: cli.Command{
			Name:      "aws-import",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return lambdaAwsImport
}

func (r LambdaAwsImport) GetCmd() cli.Command {
	return r.Command
}
