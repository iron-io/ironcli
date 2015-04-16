package httpclient

// import (
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// )
//
// type Response struct {
// 	Code    int
// 	Headers http.Header
// 	Body    []byte
// }
//
// // Do is a convenience method that (1) constructs an http.Request from the
// // given method, url, headers and body, (2) calls DoRequest(client, request, readRespBody)
// // and (3) returns the result
// func Do(client *http.Client, method string, u *url.URL, headers http.Header, body io.Reader, readRespBody bool) (*Response, error) {
// 	req, err := http.NewRequest(method, u.String(), body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers
//
// 	return DoRequest(client, req, readRespBody)
// }
//
// // DoRequest calls client.Do(req) and returns a response with the returned code
// // and headers if the request succeeded. returns the error otherwise.
// // if the request succeeded and readRespBody is true, reads and closes the entire response body
// // and puts it into response.Body. if the body read failed, returns the error.
// func DoRequest(client *http.Client, req *http.Request, readRespBody bool) (*Response, error) {
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
//
// 	out := new(Response)
// 	out.Code = resp.StatusCode
// 	out.Headers = resp.Header
// 	out.Body = nil
//
// 	if readRespBody && resp.ContentLength > 0 {
// 		bytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out.Body = bytes
// 	}
//
// 	return out, nil
// }
