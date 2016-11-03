package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/worker"
)

// TODO move this into iron_go3?
func dockerLogin(w *worker.Worker, args *map[string]string) (msg string, err error) {
	data, err := json.Marshal(args)
	reader := bytes.NewReader(data)

	req, err := http.NewRequest("POST", api.Action(w.Settings, "credentials").URL.String(), reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip/deflate")
	req.Header.Set("Authorization", "OAuth "+w.Settings.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", w.Settings.UserAgent)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if err = api.ResponseAsError(response); err != nil {
		return "", err
	}

	var res struct {
		Msg string `json:"msg"`
	}

	err = json.NewDecoder(response.Body).Decode(&res)
	return res.Msg, err
}
