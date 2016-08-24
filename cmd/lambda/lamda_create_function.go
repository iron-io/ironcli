package lambda

import (
	"fmt"

	"github.com/urfave/cli"
)

type LambdaCreateFunction struct {
	cli.Command
}

func NewLambdaCreateFunction() *LambdaCreateFunction {
	lambdaCreateFunction := &LambdaCreateFunction{
		Command: cli.Command{
			Name:      "create-function",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return lambdaCreateFunction
}

func (r LambdaCreateFunction) GetCmd() cli.Command {
	return r.Command
}
