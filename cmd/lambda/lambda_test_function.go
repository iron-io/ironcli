package lambda

import (
	"fmt"

	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaTestFunction struct {
	functionName  string
	clientContext string
	payload       string

	cli.Command
}

func NewLambdaTestFunction() *LambdaTestFunction {
	lambdaTestFunction := &LambdaTestFunction{}

	lambdaTestFunction.Command = cli.Command{
		Name:      "test-function",
		Usage:     "test function for lambda",
		ArgsUsage: "[image] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "function-name",
				Usage:       "",
				Destination: &lambdaTestFunction.functionName,
			},
			cli.StringFlag{
				Name:        "client-context",
				Usage:       "",
				Destination: &lambdaTestFunction.clientContext,
			},
			cli.StringFlag{
				Name:        "payload",
				Usage:       "",
				Destination: &lambdaTestFunction.payload,
			},
		},
		Action: func(c *cli.Context) error {
			exists, err := lambda.ImageExists(lambdaTestFunction.functionName)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("Function %s does not exist.", lambdaTestFunction.functionName)
			}

			payload := ""
			if lambdaTestFunction.payload != "" {
				payload = lambdaTestFunction.payload
			}

			// Redirect output to stdout.
			err = lambda.RunImageWithPayload(lambdaTestFunction.functionName, payload)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaTestFunction
}

func (r LambdaTestFunction) GetCmd() cli.Command {
	return r.Command
}
