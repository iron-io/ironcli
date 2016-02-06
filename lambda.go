package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/iron-io/ironcli/lambda"
)

type LambdaCreateCmd struct {
	mqCmd

	functionName *string
	runtime      *string
	handler      *string
	fileNames    []string
}

func (mf *MqFlags) functionName() *string {
	return mf.String("function-name", "", "")
}

func (mf *MqFlags) handler() *string {
	return mf.String("handler", "", "")
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
	lcc.flags = NewMqFlagSet()

	lcc.functionName = lcc.flags.functionName()
	lcc.handler = lcc.flags.handler()

	if err := lcc.flags.Parse(args); err != nil {
		return err
	}

	return lcc.flags.validateAllFlags()
}

func (lcc *LambdaCreateCmd) Run() {
	fmt.Println("fn", *lcc.functionName)

	files := make([]lambda.FileLike, 0, len(lcc.fileNames))
	for _, fileName := range lcc.fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			panic("Handle this")
		}
		files = append(files, file)
	}
	fmt.Println("OPened correctly")
	err := lambda.CreateImage(*lcc.functionName, "iron/lambda-node", *lcc.handler, files...)
	if err != nil {
		fmt.Fprintln(os.Stderr, red(err))
	}
}
