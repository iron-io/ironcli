package main

import (
	"github.com/iron-io/ironcli/cmd"
	"github.com/iron-io/ironcli/cmd/docker"
	"github.com/iron-io/ironcli/cmd/functions"
	"github.com/iron-io/ironcli/cmd/lambda"
	"github.com/iron-io/ironcli/cmd/mq"
	"github.com/iron-io/ironcli/cmd/worker"
)

func main() {
	cmd.RootCmd.AddCommand(
		functions.RootCmd,
		docker.RootCmd,
		lambda.RootCmd,
		mq.RootCmd,
		cmd.RegisterCmd,
		cmd.RunCmd,
		worker.RootCmd,
	)
	cmd.Execute()
}
