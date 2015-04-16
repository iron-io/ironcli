package httpclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/iron-io/ironcli/Godeps/_workspace/src/github.com/arschles/assert"
)

const (
	respCode = http.StatusOK
	respBody = "testbody"
	hdrKey   = "hdr1k"
	hdrVal   = "hdr1v"
)

// caller is responsible for closing the returned server
func server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(hdrKey, hdrVal)
		w.WriteHeader(respCode)
		w.Write([]byte(respBody))
	}))
}

func TestDo(t *testing.T) {
	srv := server()
	defer srv.Close()

	srvURL, err := url.Parse(srv.URL)
	assert.NoErr(t, err)

	resp, err := Do(http.DefaultClient, "GET", srvURL, nil, nil, true)
	assert.NoErr(t, err)
	assert.Equal(t, resp.Code, respCode, "response code")
	assert.Equal(t, resp.Headers.Get(hdrKey), hdrVal, "returned headers")
	assert.True(t, resp.Body != nil, "response body was nil")
	assert.Equal(t, string(resp.Body), respBody, "response body")

	resp, err = Do(http.DefaultClient, "GET", srvURL, nil, nil, false)
	assert.NoErr(t, err)
	assert.True(t, resp.Body == nil, "response body was not nil")
}

//TODO add test with no resp body
