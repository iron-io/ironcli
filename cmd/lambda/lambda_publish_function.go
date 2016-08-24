package lambda

import (
	"fmt"

	"github.com/urfave/cli"
)

type LambdaPublishFunction struct {
	cli.Command
}

func NewLambdaPublishFunction() *LambdaPublishFunction {
	lambdaPublishFunction := &LambdaPublishFunction{
		Command: cli.Command{
			Name:      "publish-function",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return lambdaPublishFunction
}

func (r LambdaPublishFunction) GetCmd() cli.Command {
	return r.Command
}
