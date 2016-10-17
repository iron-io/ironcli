package commands

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	aws_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_lambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/lambda/lambda"
)

var availableRuntimes = []string{"nodejs", "python2.7", "java8"}

const (
	skipFunctionName = iota
	requireFunctionName
)

type LambdaFlags struct {
	*flag.FlagSet
}

func (lf *LambdaFlags) validateAllFlags(fnRequired int) error {
	fn := lf.Lookup("function-name")
	// Everything except import needs a function
	if fnRequired == requireFunctionName && (fn == nil || fn.Value.String() == "") {
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
			return fmt.Errorf("Invalid runtime. Supported runtimes %s", availableRuntimes)
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

func (lf *LambdaFlags) image() *string {
	return lf.String("image", "", "By default the name of the Docker image is the name of the Lambda function. Use this to set a custom name.")
}

func (lf *LambdaFlags) version() *string {
	return lf.String("version", "$LATEST", "Version of the function to import.")
}

func (lf *LambdaFlags) downloadOnly() *bool {
	return lf.Bool("download-only", false, "Only download the function into a directory. Will not create a Docker image.")
}

func (lf *LambdaFlags) awsProfile() *string {
	return lf.String("profile", "", "AWS Profile to load from credentials file.")
}

func (lf *LambdaFlags) awsRegion() *string {
	return lf.String("region", "us-east-1", "AWS region to use.")
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

Create Docker image that can run your Lambda function. The files are the contents of the zip file to be uploaded to AWS Lambda.
`)
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

	return lcc.flags.validateAllFlags(requireFunctionName)
}

type DockerJsonWriter struct {
	under io.Writer
	w     io.Writer
}

func NewDockerJsonWriter(under io.Writer) *DockerJsonWriter {
	r, w := io.Pipe()
	go func() {
		err := jsonmessage.DisplayJSONMessagesStream(r, under, 1, true, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, red(err))
			os.Exit(1)
		}
	}()
	return &DockerJsonWriter{under, w}
}

func (djw *DockerJsonWriter) Write(p []byte) (int, error) {
	return djw.w.Write(p)
}

func (lcc *LambdaCreateCmd) Run() {
	files := make([]lambda.FileLike, 0, len(lcc.fileNames))
	opts := lambda.CreateImageOptions{
		Name:          *lcc.functionName,
		Base:          fmt.Sprintf("iron/lambda-%s", *lcc.runtime),
		Package:       "",
		Handler:       *lcc.handler,
		OutputStream:  NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	}

	if *lcc.handler == "" {
		fmt.Fprintln(os.Stderr, red("No handler specified."))
		os.Exit(1)
	}

	// For Java we allow only 1 file and it MUST be a JAR.
	if *lcc.runtime == "java8" {
		if len(lcc.fileNames) != 1 {
			fmt.Fprintln(os.Stderr, red("Java Lambda functions can only include 1 file and it must be a JAR file."))
			os.Exit(1)
		}

		if filepath.Ext(lcc.fileNames[0]) != ".jar" {
			fmt.Fprintln(os.Stderr, red("Java Lambda function package must be a JAR file."))
			os.Exit(1)
		}

		opts.Package = filepath.Base(lcc.fileNames[0])
	}

	for _, fileName := range lcc.fileNames {
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, red(err))
			os.Exit(1)
		}
		files = append(files, file)
	}

	err := lambda.CreateImage(opts, files...)
	if err != nil {
		fmt.Fprintln(os.Stderr, red(err))
		os.Exit(1)
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
	
Runs local Dockerized Lambda function and writes output to stdout.
`)
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

	return lcc.flags.validateAllFlags(requireFunctionName)
}

func (lcc *LambdaTestFunctionCmd) Run() {
	exists, err := lambda.ImageExists(*lcc.functionName)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error communicating with Docker daemon", err))
		os.Exit(1)
	}

	if !exists {
		fmt.Fprintln(os.Stderr, red(fmt.Sprintf("Function %s does not exist.", *lcc.functionName)))
		os.Exit(1)
	}

	payload := ""
	if lcc.payload != nil {
		payload = *lcc.payload
	}
	// Redirect output to stdout.
	err = lambda.RunImageWithPayload(*lcc.functionName, payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, red(err))
		os.Exit(1)
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
	
Pushes Lambda function to Docker Hub and registers with IronWorker.
If you do not want to use IronWorker, simply run 'docker push NAME' instead.
`)
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

	return lcc.flags.validateAllFlags(requireFunctionName)
}

func (lcc *LambdaPublishCmd) Run() {
	exists, err := lambda.ImageExists(*lcc.functionName)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error communicating with Docker daemon:", err))
		os.Exit(1)
	}

	if !exists {
		fmt.Fprintln(os.Stderr, red(fmt.Sprintf("Function %s does not exist:", *lcc.functionName)))
		os.Exit(1)
	}

	err = lambda.PushImage(lambda.PushImageOptions{
		NameVersion:   *lcc.functionName,
		OutputStream:  NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error pushing image:", err))
		os.Exit(1)
	}

	err = lambda.RegisterWithIron(*lcc.functionName)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error registering with IronWorker:", err))
		os.Exit(1)
	}
}

type LambdaImportCmd struct {
	lambdaCmd

	arn          string
	version      *string
	downloadOnly *bool
	awsProfile   *string
	image        *string
	awsRegion    *string
}

func (lcc *LambdaImportCmd) Args() error {
	if lcc.flags.NArg() < 1 {
		return errors.New(`import requires an AWS function ARN to import.`)
	}

	lcc.arn = lcc.flags.Arg(0)
	return nil
}

func (lcc *LambdaImportCmd) Usage() {
	fmt.Fprintln(os.Stderr, `usage: iron lambda aws-import [--region <region>] [--profile <aws profile>] [--version <version>] [--download-only] [--image <name>] ARN
	
Converts an existing Lambda function to an image. 

The function code is downloaded to a directory in the current working directory
that has the same name as the Lambda function.
`)
	lcc.flags.PrintDefaults()
}

func (lcc *LambdaImportCmd) Config() error {
	return nil
}

func (lcc *LambdaImportCmd) Flags(args ...string) error {
	flags := flag.NewFlagSet("commands", flag.ContinueOnError)
	flags.Usage = func() {}
	lcc.flags = &LambdaFlags{flags}

	lcc.version = lcc.flags.version()
	lcc.downloadOnly = lcc.flags.downloadOnly()
	lcc.awsProfile = lcc.flags.awsProfile()
	lcc.image = lcc.flags.image()
	lcc.awsRegion = lcc.flags.awsRegion()

	if err := lcc.flags.Parse(args); err != nil {
		return err
	}

	return lcc.flags.validateAllFlags(skipFunctionName)
}

func (lcc *LambdaImportCmd) downloadToFile(url string) (string, error) {
	downloadResp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer downloadResp.Body.Close()

	// zip reader needs ReaderAt, hence the indirection.
	tmpFile, err := ioutil.TempFile("", "lambda-function-")
	if err != nil {
		return "", err
	}

	io.Copy(tmpFile, downloadResp.Body)
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func (lcc *LambdaImportCmd) unzipAndGetTopLevelFiles(dst, src string) (files []lambda.FileLike, topErr error) {
	files = make([]lambda.FileLike, 0)

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer zipReader.Close()

	var fd *os.File
	for _, f := range zipReader.File {
		path := filepath.Join(dst, f.Name)
		fmt.Printf("Extracting '%s' to '%s'\n", f.Name, path)
		if f.FileInfo().IsDir() {
			os.Mkdir(path, 0644)
			// Only top-level dirs go into the list since that is what CreateImage expects.
			if filepath.Dir(f.Name) == filepath.Base(f.Name) {
				fd, topErr = os.Open(path)
				if topErr != nil {
					break
				}
				files = append(files, fd)
			}
		} else {
			// We do not close fd here since we may want to use it to dockerize.
			fd, topErr = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
			if topErr != nil {
				break
			}

			var zipFd io.ReadCloser
			zipFd, topErr = f.Open()
			if topErr != nil {
				break
			}

			_, topErr = io.Copy(fd, zipFd)
			if topErr != nil {
				// OK to skip closing fd here.
				break
			}

			zipFd.Close()

			// Only top-level files go into the list since that is what CreateImage expects.
			if filepath.Dir(f.Name) == "." {
				_, topErr = fd.Seek(0, 0)
				if topErr != nil {
					break
				}

				files = append(files, fd)
			} else {
				fd.Close()
			}
		}
	}
	return
}

func (lcc *LambdaImportCmd) getFunction() (*aws_lambda.GetFunctionOutput, error) {
	creds := aws_credentials.NewChainCredentials([]aws_credentials.Provider{
		&aws_credentials.EnvProvider{},
		&aws_credentials.SharedCredentialsProvider{
			Filename: "", // Look in default location.
			Profile:  *lcc.awsProfile,
		},
	})

	conf := aws.NewConfig().WithCredentials(creds).WithCredentialsChainVerboseErrors(true).WithRegion(*lcc.awsRegion)
	sess := aws_session.New(conf)
	conn := aws_lambda.New(sess)
	resp, err := conn.GetFunction(&aws_lambda.GetFunctionInput{
		FunctionName: aws.String(lcc.arn),
		Qualifier:    aws.String(*lcc.version),
	})

	return resp, err
}

func (lcc *LambdaImportCmd) Run() {
	function, err := lcc.getFunction()
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error getting function information", err))
		os.Exit(1)
	}
	functionName := *function.Configuration.FunctionName

	err = os.Mkdir(fmt.Sprintf("./%s", functionName), os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error creating directory: '"+functionName+"':", err))
		os.Exit(1)
	}

	tmpFileName, err := lcc.downloadToFile(*function.Code.Location)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error downloading code", err))
		os.Exit(1)
	}
	defer os.Remove(tmpFileName)

	files := make([]lambda.FileLike, 0)

	if *function.Configuration.Runtime == "java8" {
		fmt.Println("Found Java Lambda function. Going to assume code is a single JAR file.")
		path := filepath.Join(functionName, "function.jar")
		os.Rename(tmpFileName, path)
		fd, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, red(err))
			os.Exit(1)
		}

		files = append(files, fd)
	} else {
		files, err = lcc.unzipAndGetTopLevelFiles(functionName, tmpFileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, red(err))
			os.Exit(1)
		}
	}

	if *lcc.downloadOnly {
		// Since we are a command line program that will quit soon, it is OK to
		// let the OS clean `files` up.
		return
	}

	opts := lambda.CreateImageOptions{
		Name:          functionName,
		Base:          fmt.Sprintf("iron/lambda-%s", *function.Configuration.Runtime),
		Package:       "",
		Handler:       *function.Configuration.Handler,
		OutputStream:  NewDockerJsonWriter(os.Stdout),
		RawJSONStream: true,
	}

	if *lcc.image != "" {
		opts.Name = *lcc.image
	}

	if *function.Configuration.Runtime == "java8" {
		opts.Package = filepath.Base(files[0].(*os.File).Name())
	}

	err = lambda.CreateImage(opts, files...)
	if err != nil {
		fmt.Fprintln(os.Stderr, red("Error creating image", err))
		os.Exit(1)
	}
}
