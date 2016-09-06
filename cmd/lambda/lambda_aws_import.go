package lambda

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	aws_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_lambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/lambda/lambda"
	"github.com/urfave/cli"
)

type LambdaAwsImport struct {
	version      string
	downloadOnly bool
	awsProfile   string
	image        string
	awsRegion    string

	cli.Command
}

func NewLambdaAwsImport() *LambdaAwsImport {
	lambdaAwsImport := &LambdaAwsImport{}

	lambdaAwsImport.Command = cli.Command{
		Name: "aws-import",
		Usage: `Converts an existing Lambda function to an image. The function code is downloaded to a directory in the current working directory
that has the same name as the Lambda function. About ARN - (http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html).`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "version",
				Usage:       "version of the function to import.",
				Value:       "$LATEST",
				Destination: &lambdaAwsImport.version,
			},
			cli.StringFlag{
				Name:        "aws-profile",
				Usage:       "AWS Profile to load from credentials file.",
				Destination: &lambdaAwsImport.awsProfile,
			},
			cli.StringFlag{
				Name:        "image",
				Usage:       "by default the name of the Docker image is the name of the Lambda function. Use this to set a custom name.",
				Destination: &lambdaAwsImport.image,
			},
			cli.StringFlag{
				Name:        "aws-region",
				Usage:       "AWS region to use.",
				Value:       "us-east-1",
				Destination: &lambdaAwsImport.awsRegion,
			},
			cli.BoolFlag{
				Name:        "download-only",
				Usage:       "Only download the function into a directory. Will not create a Docker image.",
				Destination: &lambdaAwsImport.downloadOnly,
			},
		},
		ArgsUsage: "[ARN]",
		Action: func(c *cli.Context) error {
			function, err := lambdaAwsImport.getFunction(c.Args().First())
			if err != nil {
				return err
			}
			functionName := *function.Configuration.FunctionName

			err = os.Mkdir(fmt.Sprintf("./%s", functionName), os.ModePerm)
			if err != nil {
				return err
			}

			tmpFileName, err := lambdaAwsImport.downloadToFile(*function.Code.Location)
			if err != nil {
				return err
			}
			defer os.Remove(tmpFileName)

			files := make([]lambda.FileLike, 0)

			if *function.Configuration.Runtime == "java8" {
				fmt.Println("Found Java Lambda function. Going to assume code is a single JAR file.")
				path := filepath.Join(functionName, "function.jar")
				os.Rename(tmpFileName, path)
				fd, err := os.Open(path)
				if err != nil {
					return err
				}

				files = append(files, fd)
			} else {
				files, err = lambdaAwsImport.unzipAndGetTopLevelFiles(functionName, tmpFileName)
				if err != nil {
					return err
				}
			}

			if lambdaAwsImport.downloadOnly {
				return nil
			}

			opts := lambda.CreateImageOptions{
				Name:          functionName,
				Base:          fmt.Sprintf("iron/lambda-%s", *function.Configuration.Runtime),
				Package:       "",
				Handler:       *function.Configuration.Handler,
				OutputStream:  common.NewDockerJsonWriter(os.Stdout),
				RawJSONStream: true,
			}

			if lambdaAwsImport.image != "" {
				opts.Name = lambdaAwsImport.image
			}

			if *function.Configuration.Runtime == "java8" {
				opts.Package = filepath.Base(files[0].(*os.File).Name())
			}

			err = lambda.CreateImage(opts, files...)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return lambdaAwsImport
}

func (r LambdaAwsImport) GetCmd() cli.Command {
	return r.Command
}

func (r *LambdaAwsImport) downloadToFile(url string) (string, error) {
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

func (r *LambdaAwsImport) unzipAndGetTopLevelFiles(dst, src string) (files []lambda.FileLike, topErr error) {
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

func (r *LambdaAwsImport) getFunction(arn string) (*aws_lambda.GetFunctionOutput, error) {
	creds := aws_credentials.NewChainCredentials([]aws_credentials.Provider{
		&aws_credentials.EnvProvider{},
		&aws_credentials.SharedCredentialsProvider{
			Filename: "", // Look in default location.
			Profile:  r.awsProfile,
		},
	})

	conf := aws.NewConfig().WithCredentials(creds).WithCredentialsChainVerboseErrors(true).WithRegion(r.awsRegion)
	sess := aws_session.New(conf)
	conn := aws_lambda.New(sess)
	resp, err := conn.GetFunction(&aws_lambda.GetFunctionInput{
		FunctionName: aws.String(arn),
		Qualifier:    aws.String(r.version),
	})

	return resp, err
}
