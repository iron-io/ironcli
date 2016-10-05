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
		return fmt.Errorf(common.Red("Error communicating with Docker daemon: %v"), err)
	}

	if !exists {
		return fmt.Errorf(common.Red("Function %s does not exist."), l.FunctionName)
	}

	if !hubPushableName(l.FunctionName) {
		msg := `publish-function only supports Docker Hub right now.
		Docker Hub requires that Docker image names have the form <Docker Hub username>/<image name>:<version>. "%s" does not match this pattern.
		You may have to rename the function using "docker tag %s NEW_NAME && docker rmi %s".`
		return fmt.Errorf(msg, l.FunctionName, l.FunctionName, l.FunctionName)
	}

	err = lambda.PushImage(lambda.PushImageOptions{
		NameVersion:   l.FunctionName,
		OutputStream:  common.NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	})
	if err != nil {
		return fmt.Errorf(common.Red("Error pushing image: %v"), err)
	}

	err = lambda.RegisterWithIron(l.FunctionName)
	if err != nil {
		return fmt.Errorf(common.Red("Error registering with IronWorker: %v"), err)
	}

	return nil
}
