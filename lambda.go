package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	aws_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/lambda/lambda"
)

var availableRuntimes = []string{"nodejs", "python2.7", "java8"}

type LambdaFlags struct {
	*flag.FlagSet
}

func (lf *LambdaFlags) validateAllFlags() error {
	fn := lf.Lookup("function-name")
	if fn.Value.String() == "" {
		return errors.New(fmt.Sprintf("Please specify function-name."))
	}

	selectedRuntime := lf.Lookup("runtime")
	if selectedRuntime != nil {
		validRuntime := false
		for _, r := range availableRuntimes {
			if selectedRuntime.Value.String() == r {
				validRuntime = true
			}
		}

		if !validRuntime {
			return errors.New(fmt.Sprintf("Invalid runtime. Supported runtimes %s", availableRuntimes))
		}
	}

	return nil
}

func (lf *LambdaFlags) functionName() *string {
	return lf.String("function-name", "", "Name of function. This is usually follows Docker image naming conventions.")
}

func (lf *LambdaFlags) handler() *string {
	return lf.String("handler", "", "function/class that is the entrypoint for this function. Of the form <module name>.<function name> for nodejs/Python, <full class name>::<function name base> for Java.")
}

func (lf *LambdaFlags) runtime() *string {
	return lf.String("runtime", "", fmt.Sprintf("Runtime that your Lambda function depends on. Valid values are %s.", strings.Join(availableRuntimes, ", ")))
}

func (lf *LambdaFlags) clientContext() *string {
	return lf.String("client-context", "", "")
}

func (lf *LambdaFlags) payload() *string {
	return lf.String("payload", "", "Payload to pass to the Lambda function. This is usually a JSON object.")
}

type lambdaCmd struct {
	settings  config.Settings
	flags     *LambdaFlags
	token     *string
	projectID *string
}

type LambdaCreateCmd struct {
	lambdaCmd

	functionName *string
	runtime      *string
	handler      *string
	fileNames    []string
}

func (lcc *LambdaCreateCmd) Args() error {
	if lcc.flags.NArg() < 1 {
		return errors.New(`lambda create requires at least one file`)
	}

	for _, arg := range lcc.flags.Args() {
		lcc.fileNames = append(lcc.fileNames, arg)
	}

	return nil
}

func (lcc *LambdaCreateCmd) Usage() {
	fmt.Fprintln(os.Stderr, `usage: iron lambda create-function --function-name NAME --runtime RUNTIME --handler HANDLER file [files...]

Create Docker image that can run your Lambda function. The files are the contents of the zip file to be uploaded to AWS Lambda.`)
	lcc.flags.PrintDefaults()
}

func (lcc *LambdaCreateCmd) Config() error {
	return nil
}

func (lcc *LambdaCreateCmd) Flags(args ...string) error {
	flags := flag.NewFlagSet("commands", flag.ContinueOnError)
	flags.Usage = func() {}
	lcc.flags = &LambdaFlags{flags}

	lcc.functionName = lcc.flags.functionName()
	lcc.handler = lcc.flags.handler()
	lcc.runtime = lcc.flags.runtime()

	if err := lcc.flags.Parse(args); err != nil {
		return err
	}

	return lcc.flags.validateAllFlags()
}

type DockerJsonWriter struct {
	under io.Writer
	w     io.Writer
}

func NewDockerJsonWriter(under io.Writer) *DockerJsonWriter {
	r, w := io.Pipe()
	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(r, os.Stdout, 1, true, nil)
		log.Fatal(err)
	}()
	return &DockerJsonWriter{under, w}
}

func (djw *DockerJsonWriter) Write(p []byte) (int, error) {
	return djw.w.Write(p)
}

func (lcc *LambdaCreateCmd) Run() {
	files := make([]lambda.FileLike, 0, len(lcc.fileNames))
	for _, fileName := range lcc.fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, file)
	}
	err := lambda.CreateImage(lambda.CreateImageOptions{*lcc.functionName, fmt.Sprintf("iron/lambda-%s", *lcc.runtime), *lcc.handler, NewDockerJsonWriter(os.Stdout)}, files...)
	if err != nil {
		log.Fatal(err)
	}
}

type LambdaTestFunctionCmd struct {
	lambdaCmd

	functionName  *string
	clientContext *string
	payload       *string
}

func (lcc *LambdaTestFunctionCmd) Args() error {
	return nil
}

func (lcc *LambdaTestFunctionCmd) Usage() {
	fmt.Fprintln(os.Stderr, `usage: iron lambda test-function --function-name NAME [--client-context <value>] [--payload <value>]
	
Runs local Dockerized Lambda function and writes output to stdout.`)
	lcc.flags.PrintDefaults()
}

func (lcc *LambdaTestFunctionCmd) Config() error {
	return nil
}

func (lcc *LambdaTestFunctionCmd) Flags(args ...string) error {
	flags := flag.NewFlagSet("commands", flag.ContinueOnError)
	flags.Usage = func() {}
	lcc.flags = &LambdaFlags{flags}

	lcc.functionName = lcc.flags.functionName()
	lcc.clientContext = lcc.flags.clientContext()
	lcc.payload = lcc.flags.payload()

	if err := lcc.flags.Parse(args); err != nil {
		return err
	}

	return lcc.flags.validateAllFlags()
}

func (lcc *LambdaTestFunctionCmd) Run() {
	exists, err := lambda.ImageExists(*lcc.functionName)
	if err != nil {
		log.Fatal("Error communicating with Docker daemon", err)
	}

	if !exists {
		log.Fatal(fmt.Sprintf("Function %s does not exist.", *lcc.functionName))
	}

	payload := ""
	if lcc.payload != nil {
		payload = *lcc.payload
	}
	// Redirect output to stdout.
	err = lambda.RunImageWithPayload(*lcc.functionName, payload)
	if err != nil {
		log.Fatal(err)
	}
}

type LambdaPublishCmd struct {
	lambdaCmd

	functionName *string
}

func (lcc *LambdaPublishCmd) Args() error {
	return nil
}

func (lcc *LambdaPublishCmd) Usage() {
	fmt.Fprintln(os.Stderr, `usage: iron lambda publish-function --function-name NAME
	
Pushes Lambda function to Docker Hub and registers with IronWorker.`)
	lcc.flags.PrintDefaults()
}

func (lcc *LambdaPublishCmd) Config() error {
	return nil
}

func (lcc *LambdaPublishCmd) Flags(args ...string) error {
	flags := flag.NewFlagSet("commands", flag.ContinueOnError)
	flags.Usage = func() {}
	lcc.flags = &LambdaFlags{flags}

	lcc.functionName = lcc.flags.functionName()

	if err := lcc.flags.Parse(args); err != nil {
		return err
	}

	return lcc.flags.validateAllFlags()
}

func (lcc *LambdaPublishCmd) Run() {
	exists, err := lambda.ImageExists(*lcc.functionName)
	if err != nil {
		log.Fatal("Error communicating with Docker daemon:", err)
	}

	if !exists {
		log.Fatal(fmt.Sprintf("Function %s does not exist:", *lcc.functionName))
	}

	err = lambda.PushImage(*lcc.functionName)
	if err != nil {
		log.Fatal("Error pushing image:", err)
	}

	err = lambda.RegisterWithIron(*lcc.functionName, aws_credentials.NewEnvCredentials())
	if err != nil {
		log.Fatal("Error registering with IronWorker:", err)
	}
}
