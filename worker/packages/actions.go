package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/iron-io/ironcli/common"
	"github.com/iron-io/ironcli/httpclient"
)

func baseReq(g *common.GlobalFlags, method string, body string, pathFmt string, vals ...interface{}) *http.Request {
	pathStr := fmt.Sprintf(pathFmt, vals)
	urlStr := fmt.Sprintf("https://%s.iron.io/%d/projects/%s/%s", g.Host, g.Version, g.ProjID, pathStr)
	req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(body))
	if err != nil {
		// TODO something smarter here
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", g.Token))
	return req
}

func decode(resp *httpclient.Response, i interface{}) {
	err := json.Unmarshal(resp.Body, i)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doJSON(req *http.Request, resp interface{}) {
	rawResp, err := httpclient.DoRequest(http.DefaultClient, req, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decode(rawResp, resp)
}

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
	req := baseReq(g, "GET", "", "codes")

	var resp struct {
		Codes []code `json:"codes"`
	}
	doJSON(req, &resp)
	common.PrintJSON(resp)
}

func upload(g *common.GlobalFlags) {
	fmt.Println("TODO")
}

func info(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := baseReq(g, "GET", "", "codes/%s", codeID)
	resp := code{}
	doJSON(req, &resp)
	common.PrintJSON(resp)
}

func del(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := baseReq(g, "DELETE", "", "codes/%s", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	doJSON(req, &resp)
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
	req := baseReq(g, "GET", "", "codes/%s?page=%d&per_page=%d", codeID, page, perPage)
	var resp struct {
		Revs []rev `json:"revisions"`
	}
	doJSON(req, &resp)
	common.PrintJSON(resp)
}

func pause(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := baseReq(g, "POST", "", "codes/%s/pause_task_queue", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	doJSON(req, &resp)
	common.PrintJSON(resp)
}

func resume(g *common.GlobalFlags) {
	codeID := g.Ctx.String("codeid")
	req := baseReq(g, "POST", "", "codes/%s/resume_task_queue", codeID)
	var resp struct {
		Msg string `json:"msg"`
	}
	doJSON(req, &resp)
	common.PrintJSON(resp)
}
