package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/ironcli/common"
	"github.com/urfave/cli"
)

type DockerLogin struct {
	Email      string
	Password   string
	Url        string
	Username   string
	TestAuth   string `json:"-"`
	RemoteAuth string `json:"-"`
	Auth       string `json:"auth"`

	cli.Command
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
				Usage:       "docker repo url, if you're using custom repo",
				Destination: &dockerLogin.Url,
			},
			cli.StringFlag{
				Name:        "username",
				Usage:       "docker repo user name",
				Destination: &dockerLogin.Username,
			},
		},
		Action: func(c *cli.Context) error {
			err := dockerLogin.login()
			if err != nil {
				return err
			}

			auth := map[string]string{
				"auth": dockerLogin.RemoteAuth,
			}

			msg, err := dockerLogin.Execute(&settings.Worker, &auth)
			if err != nil {
				return err
			}

			fmt.Println(`Added docker repo credentials: ` + msg)

			return nil
		},
	}

	return dockerLogin
}

func (r DockerLogin) GetCmd() cli.Command {
	return r.Command
}

func (r *DockerLogin) login() error {
	if r.Url == "" {
		defaultUrl := "https://index.docker.io/v1/"
		r.Url = defaultUrl
	}

	auth := base64.StdEncoding.EncodeToString([]byte(r.Username + ":" + r.Password))
	r.TestAuth = auth

	bytes, _ := json.Marshal(*r)
	authString := base64.StdEncoding.EncodeToString(bytes)
	r.RemoteAuth = authString

	req, err := http.NewRequest("GET", r.Url+"users/", nil)
	if err != nil {
		return fmt.Errorf("error authenticating docker login: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+r.TestAuth)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return errors.New("Docker repo auth failed")
	}

	return nil
}

func (r *DockerLogin) Execute(settings *config.Settings, args *map[string]string) (msg string, err error) {
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
