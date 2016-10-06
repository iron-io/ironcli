package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/iron-io/iron_go3/config"
)

type FunctionFlags struct {
	*flag.FlagSet
}

func (ff *FunctionFlags) validateAllFlags() error {
	// This is how you validate individual flags as having valid values.
	fn := ff.Lookup("image")

	// Just an example!
	if fn == nil || len(fn.Value.String()) < 3 {
		return errors.New(fmt.Sprintf("Function name has to be greater than 3 chars."))
	}

	// All flags good!
	return nil
}

func (ff *FunctionFlags) image() *string {
	return ff.String("image", "thisisadefaultvalue", "Image to call when endpoint is hit.")
}

func (ff *FunctionFlags) path() *string {
	return ff.String("path", "", "specify route")
}

type functionCmd struct {
	settings  config.Settings
	flags     *FunctionFlags
	token     *string
	projectID *string
}

type FunctionCreateCmd struct {
	functionCmd

	image *string
	path  *string
}

func (fcc *FunctionCreateCmd) Args() error {
	// Use this to process any positional arguments like list of files etc.
	//for _, arg := range fcc.flags.Args() {
	// do something...
	//}

	return nil
}

func (fcc *FunctionCreateCmd) Usage() {
	fmt.Fprintln(os.Stderr, `usage: iron function create-function --path PATH --image IMAGE

Put description here.
`)
	fcc.flags.PrintDefaults()
}

func (fcc *FunctionCreateCmd) Config() error {
	// Can be used to load configuration settings. Should be unused since functions is not tied to auth or anything. See mq.go if you want a reference.
	return nil
}

// This is where we actually initialize the flags and assign them to variables. All of this is documented in the golang flag package.
func (fcc *FunctionCreateCmd) Flags(args ...string) error {
	// setup flags.
	flags := flag.NewFlagSet("commands", flag.ContinueOnError)
	flags.Usage = func() {}
	fcc.flags = &FunctionFlags{flags}

	fcc.path = fcc.flags.path()
	fcc.image = fcc.flags.image()

	// parse actual args.
	if err := fcc.flags.Parse(args); err != nil {
		return err
	}

	// validate them
	return fcc.flags.validateAllFlags()
}

func (fcc *FunctionCreateCmd) Run() {
	// Boom! Use the parameters as you'd like.
	if fcc.path == nil {
		// argument was never passed.
		fmt.Fprintln(os.Stderr, red("Path required"))
	}

	if *fcc.path == "" {
		// argument was passed but empty.
		fmt.Fprintln(os.Stderr, red("Path required"))
	}

	// This is made up.
	//functions.CreateFunction(*fcc.path, *fcc.image)
}
