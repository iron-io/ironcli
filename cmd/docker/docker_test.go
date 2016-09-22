package docker

import (
	"os"
	"testing"

	"github.com/iron-io/ironcli/common"
)

func TestDockerLogin(t *testing.T) {
	var (
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)

	dockerLogin := NewDockerLogin(settings)
	dockerLogin.Email = os.Getenv("TEST_DOCKER_EMAIL")
	dockerLogin.Username = os.Getenv("TEST_DOCKER_USERNAME")
	dockerLogin.Password = os.Getenv("TEST_DOCKER_PASSWORD")

	err := dockerLogin.Action(settings)
	if err != nil {
		t.Error(err)
	}
}
