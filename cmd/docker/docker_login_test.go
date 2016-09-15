package docker

import (
	"flag"
	"testing"

	"github.com/iron-io/ironcli/common"
	"github.com/stretchr/testify/assert"
)

var (
	email    string
	username string
	password string
)

func init() {
	flag.StringVar(&email, "email", "", "")
	flag.StringVar(&username, "username", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.Parse()
}

func TestDockerLogin(t *testing.T) {
	var (
		assert   = assert.New(t)
		settings = &common.Settings{Product: "iron_worker"}
	)

	common.SetSettings(settings)
	dockerLogin := NewDockerLogin(settings)
	dockerLogin.Email = email
	dockerLogin.Username = username
	dockerLogin.Password = password

	err := dockerLogin.login()
	assert.Nil(err)

	auth := map[string]string{
		"auth": dockerLogin.RemoteAuth,
	}

	_, err = dockerLogin.Execute(&settings.Worker, &auth)
	assert.Nil(err)
}
