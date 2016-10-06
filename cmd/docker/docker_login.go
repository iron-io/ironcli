package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type DockerLogin struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ServerAddress string `json:"serveraddress"`

	cli.Command `json:"-"`
}

func NewDockerLogin(settings *common.Settings) *DockerLogin {
	dockerLogin := &DockerLogin{}
	dockerLogin.Command = cli.Command{
		Name:      "login",
		Usage:     "manage docker registry credentials.",
		ArgsUsage: "[args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "email",
				Usage:       "docker repo user email",
				Destination: &dockerLogin.Email,
			},
			cli.StringFlag{
				Name:        "password",
				Usage:       "docker repo password",
				Destination: &dockerLogin.Password,
			},
			cli.StringFlag{
				Name:        "url",
				Usage:       "docker repo url, if you're using custom repo, e.g. https://registry.hub.docker.com",
				Destination: &dockerLogin.ServerAddress,
			},
			cli.StringFlag{
				Name:        "username",
				Usage:       "docker repo user name",
				Destination: &dockerLogin.Username,
			},
		},
		Before: func(c *cli.Context) error {
			if err := common.SetSettings(settings); err != nil {
				return err
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			err := dockerLogin.Action(settings)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return dockerLogin
}

func (d DockerLogin) GetCmd() cli.Command {
	return d.Command
}

func (d *DockerLogin) Execute(settings *config.Settings, args *map[string]string) (msg string, err error) {
	data, err := json.Marshal(args)
	reader := bytes.NewReader(data)

	req, err := http.NewRequest("POST", api.Action(*settings, "credentials").URL.String(), reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "OAuth "+settings.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", settings.UserAgent)

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

func (d *DockerLogin) Action(settings *common.Settings) error {
	bytes, err := json.Marshal(*d)
	if err != nil {
		return err
	}
	authString := base64.StdEncoding.EncodeToString(bytes)

	auth := map[string]string{
		"auth": authString,
	}

	msg, err := d.Execute(&settings.Worker, &auth)
	if err != nil {
		return err
	}

	fmt.Println(`Added docker repo credentials: ` + msg)

	return nil
}
