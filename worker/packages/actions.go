package packages

import (
	"fmt"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/iron-io/iron_go3/worker"
	"github.com/iron-io/ironcli/common"
)

func list(g *common.GlobalFlags) {
	page := g.Ctx.Int("page")
	perPage := g.Ctx.Int("per-page")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	codes, err := wrk.CodePackageList(page, perPage)
	common.FailErr(err)
	common.PrintJSON(codes)
}

func upload(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func info(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	code, err := wrk.CodePackageInfo(codeID)
	common.FailErr(err)
	common.PrintJSON(code)
}

func del(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	err := wrk.CodePackageDelete(codeID)
	common.FailErr(err)
	fmt.Println("success")
}

func download(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func listrevs(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	code, err := wrk.CodePackageRevisions(codeID)
	common.FailErr(err)
	common.PrintJSON(code)
}
