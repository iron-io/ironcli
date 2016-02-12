package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/lambda"
)

type LambdaFlags struct {
	*flag.FlagSet
}

func (lf *LambdaFlags) validateAllFlags() error {
	availableRuntimes := []string{"nodejs", "python2.7", "java8"}
	selectedRuntime := lf.Lookup("runtime")
	validRuntime := false
	for _, r := range availableRuntimes {
		if selectedRuntime.Value.String() == r {
			validRuntime = true
		}
	}

	if !validRuntime {
		return errors.New(fmt.Sprintf("Invalid runtime. Supported runtimes %s", availableRuntimes))
	}

	return nil
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

func (lf *LambdaFlags) functionName() *string {
	return lf.String("function-name", "", "")
}

func (lf *LambdaFlags) handler() *string {
	return lf.String("handler", "", "")
}

func (lf *LambdaFlags) runtime() *string {
	return lf.String("runtime", "", "")
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
	fmt.Fprintln(os.Stderr, `usage: iron lambda create-function --function-name NAME --runtime RUNTIME --handler HANDLER file [files...]`)
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

func (lcc *LambdaCreateCmd) Run() {
	files := make([]lambda.FileLike, 0, len(lcc.fileNames))
	for _, fileName := range lcc.fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, file)
	}
	err := lambda.CreateImage(*lcc.functionName, "iron/lambda-nodejs", *lcc.handler, files...)
	if err != nil {
		log.Fatal(err)
	}
}
