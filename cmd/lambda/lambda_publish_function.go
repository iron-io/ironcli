package lambda

import (
	"fmt"
	"os"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaPublishFunction struct {
	functionName string

	cli.Command
}

func NewLambdaPublishFunction() *LambdaPublishFunction {
	lambdaPublishFunction := &LambdaPublishFunction{}

	lambdaPublishFunction.Command = cli.Command{
		Name:  "publish-function",
		Usage: "publish function for lambda",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "function-name",
				Usage:       "function-name usage",
				Destination: &lambdaPublishFunction.functionName,
			},
		},
		ArgsUsage: "[NAME] [args]",
		Action: func(c *cli.Context) error {
			exists, err := lambda.ImageExists(lambdaPublishFunction.functionName)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("Function %s does not exist:", lambdaPublishFunction.functionName)
			}

			err = lambda.PushImage(lambda.PushImageOptions{
				NameVersion:   lambdaPublishFunction.functionName,
				OutputStream:  common.NewDockerJsonWriter(os.Stdout),
				RawJSONStream: true,
			})
			if err != nil {
				return err
			}

			err = lambda.RegisterWithIron(lambdaPublishFunction.functionName)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaPublishFunction
}

func (r LambdaPublishFunction) GetCmd() cli.Command {
	return r.Command
}
