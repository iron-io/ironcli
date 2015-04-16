package common

import (
	"fmt"
	"os"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/codegangsta/cli"
)

type GlobalFlags struct {
	ProjID  string
	Token   string
	Version int
	Host    string
	Ctx     *cli.Context
}

func GetGlobalFlags(c *cli.Context) *GlobalFlags {
	projID := c.GlobalString(ProjectID)
	token := c.GlobalString(Token)
	vsn := c.GlobalInt(Version)
	host := c.GlobalString(Host)
	if projID == "" {
		fmt.Println("no project ID specified")
		os.Exit(1)
	}
	if token == "" {
		fmt.Println("no token specified")
		os.Exit(1)
	}
	if vsn == InvalidVersion {
		fmt.Println("no version specified")
		os.Exit(1)
	}
	if host == "" {
		fmt.Println("no host specified")
		os.Exit(1)
	}

	return &GlobalFlags{ProjID: projID, Token: token, Version: vsn, Host: host, Ctx: c}
}

func WithGlobalFlags(fn func(g *GlobalFlags)) func(*cli.Context) {
	return func(c *cli.Context) {
		gflags := GetGlobalFlags(c)
		fn(gflags)
	}
}

// IntOrFail returns g.Ctx.Int(name). if the value is equal to missing,
// prints an error message and calls os.Exit(1)
func (g *GlobalFlags) IntOrFail(name string, missing int) int {
	i := g.Ctx.Int(name)
	if i == missing {
		fmt.Println("no ", name, " specified")
		os.Exit(1)
	}
	return i
}

// StringOrFail returns g.Ctx.String(name). if the value is missing,
// prints an error message and call os.Exit(1)
func (g *GlobalFlags) StringOrFail(name string) string {
	s := g.Ctx.String(name)
	if s == "" {
		fmt.Println("no ", name, " specified")
		os.Exit(1)
	}
	return s
}
