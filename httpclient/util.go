package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/iron-io/go/pushd/httpclient"
	"github.com/iron-io/ironcli/common"
)

// BaseReq makes a new http.Request for talking to Iron's APIs. the request will
// have the given body and Authorization: Oauth t.Token in its header. finally,
// the URL of the request will be constructed as
// https://{g.Host}.iron.io/{g.Version}/projects/{g.ProjID}/. the remainder of the URL
// will be constructed from calling fmt.Sprintf(pathFmt, vals...).
// BaseReq will print an error to stdout and os.Exit(1) if there was an error
// creating the request.
func BaseReq(g *common.GlobalFlags, method string, body string, pathFmt string, vals ...interface{}) *http.Request {
	pathStr := fmt.Sprintf(pathFmt, vals...)
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

// DecodeJSON decodes resp.Body into i, just as json.Unmarshal would.
// prints an error to stdout and calls os.Exit(1) if there was an error
// unmarshaling.
func DecodeJSON(resp *httpclient.Response, i interface{}) {
	err := json.Unmarshal(resp.Body, i)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// DoJSON is a convenience function for calling httpclient.DoRequest
// with req and then decoding the response body with json.Unmarshal
// into resp. example usage:
//  req := BaseReq(g, "GET", "", "a/b/c")
//  var resp struct{
//    A string
//  }
//  DoJSON(req, &resp)
// DoJSON will print to stdout and os.Exit(1) if there was a failure
// creating the request or it could not decode the JSON
func DoJSON(req *http.Request, resp interface{}) {
	rawResp, err := httpclient.DoRequest(http.DefaultClient, req, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	DecodeJSON(rawResp, resp)
}
