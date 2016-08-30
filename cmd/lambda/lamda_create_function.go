package lambda

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaCreateFunction struct {
	functionName string
	runtime      string
	handler      string
	fileNames    []string

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
				Usage:       "function-name usage",
				Destination: &lambdaCreateFunction.functionName,
			},
			cli.StringFlag{
				Name:        "runtime",
				Usage:       "runtime usage",
				Destination: &lambdaCreateFunction.runtime,
			},
			cli.StringFlag{
				Name:        "handler",
				Usage:       "handler usage",
				Destination: &lambdaCreateFunction.handler,
			},
		},
		Action: func(c *cli.Context) error {
			lambdaCreateFunction.fileNames = append(lambdaCreateFunction.fileNames, c.Args().First())
			lambdaCreateFunction.fileNames = append(lambdaCreateFunction.fileNames, c.Args().Tail()...)

			files := make([]lambda.FileLike, 0, len(lambdaCreateFunction.fileNames))
			opts := lambda.CreateImageOptions{
				Name:          lambdaCreateFunction.functionName,
				Base:          fmt.Sprintf("iron/lambda-%s", lambdaCreateFunction.runtime),
				Package:       "",
				Handler:       lambdaCreateFunction.handler,
				OutputStream:  common.NewDockerJsonWriter(os.Stdout),
				RawJSONStream: true,
			}

			if lambdaCreateFunction.handler == "" {
				return errors.New("No handler specified.")
			}

			// For Java we allow only 1 file and it MUST be a JAR.
			if lambdaCreateFunction.runtime == "java8" {
				if len(lambdaCreateFunction.fileNames) != 1 {
					return errors.New("Java Lambda functions can only include 1 file and it must be a JAR file.")
				}

				if filepath.Ext(lambdaCreateFunction.fileNames[0]) != ".jar" {
					return errors.New("Java Lambda function package must be a JAR file.")
				}

				opts.Package = filepath.Base(lambdaCreateFunction.fileNames[0])
			}

			for _, fileName := range lambdaCreateFunction.fileNames {
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
		},
	}

	return lambdaCreateFunction
}

func (r LambdaCreateFunction) GetCmd() cli.Command {
	return r.Command
}
