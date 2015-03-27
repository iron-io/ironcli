package stacks

import (
	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/ironcli/httpclient"
)

var Subcommands = []cli.Command{
	cli.Command{
		Name:  "list",
		Usage: "list all available stacks",
		Action: common.WithGlobalFlags(func(g *common.GlobalFlags) {
			req := httpclient.BaseReq(g, "GET", "", "stacks")
			var resp []string
			httpclient.DoJSON(req, &resp)
			common.PrintJSON(resp)
		}),
	},
}
