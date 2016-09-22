package lambda

import (
	"fmt"
	"os"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaPublishFunction struct {
	FunctionName string

	cli.Command
}

func NewLambdaPublishFunction() *LambdaPublishFunction {
	lambdaPublishFunction := &LambdaPublishFunction{}

	lambdaPublishFunction.Command = cli.Command{
		Name:  "publish-function",
		Usage: "pushes Lambda function to Docker Hub and registers with IronWorker. If you do not want to use IronWorker, simply run 'docker push NAME' instead.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "function-name",
				Usage:       "name of function. This is usually follows Docker image naming conventions.",
				Destination: &lambdaPublishFunction.FunctionName,
			},
		},
		ArgsUsage: "",
		Action: func(c *cli.Context) error {
			err := lambdaPublishFunction.Action()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaPublishFunction
}

func (l LambdaPublishFunction) GetCmd() cli.Command {
	return l.Command
}

func (l *LambdaPublishFunction) Action() error {
	exists, err := lambda.ImageExists(l.FunctionName)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Function %s does not exist:", l.FunctionName)
	}

	err = lambda.PushImage(lambda.PushImageOptions{
		NameVersion:   l.FunctionName,
		OutputStream:  common.NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	})
	if err != nil {
		return err
	}

	err = lambda.RegisterWithIron(l.FunctionName)
	if err != nil {
		return err
	}

	return nil
}
