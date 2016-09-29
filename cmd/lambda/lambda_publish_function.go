package lambda

import (
	"fmt"
	"os"
	"strings"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

func hubPushableName(image string) bool {
	return strings.Count(image, "/") == 1
}

type LambdaPublishFunction struct {
	functionName string

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

			if !hubPushableName(lambdaPublishFunction.functionName) {
				msg := `publish-function only supports Docker Hub right now.
				Docker Hub requires that Docker image names have the form <Docker Hub username>/<image name>:<version>. "%s" does not match this pattern.
				You may have to rename the function using "docker tag %s NEW_NAME && docker rmi %s".`
				return fmt.Errorf(msg, lambdaPublishFunction.functionName, lambdaPublishFunction.functionName, lambdaPublishFunction.functionName)
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
