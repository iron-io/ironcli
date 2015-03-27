package packages

import (
	"fmt"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/ironcli/httpclient"
)

type code struct {
	ID              string `json:"id"`
	ProjID          string `json:"project_id"`
	Name            string `json:"name"`
	Runtime         string `json:"runtime"`
	LatestChecksum  string `json:"latest_checksum"`
	Rev             int    `json:"rev"`
	LatestHistoryID string `json:"latest_history_id"`
	LatestChange    int    `json:"latest_change"`
}

func list(g *common.GlobalFlags) {
	req := httpclient.BaseReq(g, "GET", "", "codes")

	var resp struct {
		Codes []code `json:"codes"`
	}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}

func upload(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func info(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := httpclient.BaseReq(g, "GET", "", "codes/%s", codeID)
	resp := code{}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}

func del(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := httpclient.BaseReq(g, "DELETE", "", "codes/%s", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}

func download(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

type rev struct {
	ID       string `json:"id"`
	CodeID   string `json:"code_id"`
	ProjID   string `json:"project_id"`
	Rev      int    `json:"rev"`
	Runtime  string `json:"runtime"`
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

func listrevs(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	page := g.Ctx.Int("page")
	perPage := g.Ctx.Int("perpage")
	req := httpclient.BaseReq(g, "GET", "", "codes/%s?page=%d&per_page=%d", codeID, page, perPage)
	var resp struct {
		Revs []rev `json:"revisions"`
	}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}

func pause(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := httpclient.BaseReq(g, "POST", "", "codes/%s/pause_task_queue", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}

func resume(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := httpclient.BaseReq(g, "POST", "", "codes/%s/resume_task_queue", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	httpclient.DoJSON(req, &resp)
	common.PrintJSON(resp)
}
