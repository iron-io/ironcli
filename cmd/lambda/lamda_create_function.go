package lambda

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

var availableRuntimes = []string{"nodejs", "python2.7", "java8"}

type LambdaCreateFunction struct {
	FunctionName string
	Runtime      string
	Handler      string
	FileNames    []string

	cli.Command
}

func NewLambdaCreateFunction() *LambdaCreateFunction {
	lambdaCreateFunction := &LambdaCreateFunction{}

	lambdaCreateFunction.Command = cli.Command{
		Name:      "create-function",
		Usage:     "create function for lambda",
		ArgsUsage: "[file] [files, ...] [args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "function-name",
				Usage:       "name of function. This is usually follows Docker image naming conventions.",
				Destination: &lambdaCreateFunction.FunctionName,
			},
			cli.StringFlag{
				Name:        "runtime",
				Usage:       fmt.Sprintf("Runtime that your Lambda function depends on. Valid values are %s.", strings.Join(availableRuntimes, ", ")),
				Destination: &lambdaCreateFunction.Runtime,
			},
			cli.StringFlag{
				Name:        "handler",
				Usage:       "function/class that is the entrypoint for this function. Of the form <module name>.<function name> for nodejs/Python, <full class name>::<function name base> for Java.",
				Destination: &lambdaCreateFunction.Handler,
			},
		},
		Action: func(c *cli.Context) error {
			lambdaCreateFunction.FileNames = append(lambdaCreateFunction.FileNames, c.Args().First())
			lambdaCreateFunction.FileNames = append(lambdaCreateFunction.FileNames, c.Args().Tail()...)

			err := lambdaCreateFunction.Action()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaCreateFunction
}

func (l LambdaCreateFunction) GetCmd() cli.Command {
	return l.Command
}

func (l *LambdaCreateFunction) Action() error {
	files := make([]lambda.FileLike, 0, len(l.FileNames))
	opts := lambda.CreateImageOptions{
		Name:          l.FunctionName,
		Base:          fmt.Sprintf("iron/lambda-%s", l.Runtime),
		Package:       "",
		Handler:       l.Handler,
		OutputStream:  common.NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	}

	if l.Handler == "" {
		return errors.New("No handler specified.")
	}

	// For Java we allow only 1 file and it MUST be a JAR.
	if l.Runtime == "java8" {
		if len(l.FileNames) != 1 {
			return errors.New("Java Lambda functions can only include 1 file and it must be a JAR file.")
		}

		if filepath.Ext(l.FileNames[0]) != ".jar" {
			return errors.New("Java Lambda function package must be a JAR file.")
		}

		opts.Package = filepath.Base(l.FileNames[0])
	}

	for _, fileName := range l.FileNames {
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			return err
		}

		files = append(files, file)
	}

	err := lambda.CreateImage(opts, files...)
	if err != nil {
		return err
	}

	return nil
}
