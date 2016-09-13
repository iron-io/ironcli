package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func projectName(settings *Settings) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s://%s:%d/%s/projects/%s?oauth=%s",
		settings.Worker.Scheme, settings.Worker.Host, settings.Worker.Port,
		settings.Worker.ApiVersion, settings.Worker.ProjectId, settings.Worker.Token))

	if err != nil {
		return "", err
	}

	var reply struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)

	return reply.Name, err
}
