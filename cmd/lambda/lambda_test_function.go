package lambda

import (
	"fmt"

	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaTestFunction struct {
	FunctionName  string
	ClientContext string
	Payload       string

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
				Usage:       "name of function. This is usually follows Docker image naming conventions.",
				Destination: &lambdaTestFunction.FunctionName,
			},
			cli.StringFlag{
				Name:        "client-context",
				Destination: &lambdaTestFunction.ClientContext,
			},
			cli.StringFlag{
				Name:        "Payload",
				Usage:       "give function Payload",
				Destination: &lambdaTestFunction.Payload,
			},
		},
		Action: func(c *cli.Context) error {
			err := lambdaTestFunction.Action()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaTestFunction
}

func (l LambdaTestFunction) GetCmd() cli.Command {
	return l.Command
}

func (l *LambdaTestFunction) Action() error {
	exists, err := lambda.ImageExists(l.FunctionName)
	if err != nil {
		return fmt.Errorf("Error communicating with Docker daemon: %v", err)
	}

	if !exists {
		return fmt.Errorf("Function %s does not exist.", l.FunctionName)
	}

	payload := ""
	if l.Payload != "" {
		payload = l.Payload
	}

	// Redirect output to stdout.
	err = lambda.RunImageWithPayload(l.FunctionName, payload)
	if err != nil {
		return err
	}

	return nil
}
