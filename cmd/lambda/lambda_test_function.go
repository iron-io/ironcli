package lambda

import (
	"fmt"

	"github.com/urfave/cli"
)

type LambdaTestFunction struct {
	cli.Command
}

func NewLambdaTestFunction() *LambdaTestFunction {
	lambdaTestFunction := &LambdaTestFunction{
		Command: cli.Command{
			Name:      "test-function",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[image] [args]",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: test ", c.Args().First())
				return nil
			},
		},
	}

	return lambdaTestFunction
}

func (r LambdaTestFunction) GetCmd() cli.Command {
	return r.Command
}
