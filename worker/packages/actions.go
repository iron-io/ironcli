package packages

import (
	"fmt"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/iron-io/iron_go3/worker"
	"github.com/iron-io/ironcli/common"
)

func list(g *common.GlobalFlags) {
	page := g.IntOrFail("page", -1)
	perPage := g.IntOrFail("per-page", -1)
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	codes, err := wrk.CodePackageList(page, perPage)
	common.FailErr(err)
	common.PrintJSON(codes)
}

func upload(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func info(g *common.GlobalFlags) {
	codeID := g.StringOrFail("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	code, err := wrk.CodePackageInfo(codeID)
	common.FailErr(err)
	common.PrintJSON(code)
}

func del(g *common.GlobalFlags) {
	codeID := g.StringOrFail("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	err := wrk.CodePackageDelete(codeID)
	common.FailErr(err)
	fmt.Println("success")
}

func download(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func listrevs(g *common.GlobalFlags) {
	codeID := g.StringOrFail("codeid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	code, err := wrk.CodePackageRevisions(codeID)
	common.FailErr(err)
	common.PrintJSON(code)
}

func logs(g *common.GlobalFlags) {
	taskID := g.StringOrFail("taskid")
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	log, err := wrk.TaskLog(taskID)
	common.FailErr(err)
	fmt.Println(string(log))
}

func schedule(g *common.GlobalFlags) {
	sched := getSchedule(g)
	wrk := worker.Worker{Settings: *common.NewIronConfig(g)}
	ids, err := wrk.Schedule(sched)
	common.FailErr(err)
	common.PrintJSON(ids)
}
