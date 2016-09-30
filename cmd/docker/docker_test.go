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
	if dockerLogin.Email == "" {
		t.Error("Email shouldn't be an empty")
	}

	dockerLogin.Username = os.Getenv("TEST_DOCKER_USERNAME")
	if dockerLogin.Username == "" {
		t.Error("Username shouldn't be an empty")
	}

	dockerLogin.Password = os.Getenv("TEST_DOCKER_PASSWORD")
	if dockerLogin.Password == "" {
		t.Error("Password shouldn't be an empty")
	}

	err := dockerLogin.Action(settings)
	if err != nil {
		t.Error(err)
	}
}
