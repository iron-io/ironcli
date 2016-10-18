package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/giantswarm/semver-bump/bump"
	"github.com/giantswarm/semver-bump/storage"
	funcs "github.com/iron-io/functions/api/models"

	"gopkg.in/yaml.v2"
)

var (
	ErrLoadFunctionConfig = errors.New("Failed to load function's config")
	ErrInvalidVersion     = errors.New("Function's VERSION could not be read")
	ErrBuildFunctionImage = errors.New("Couldn't build the function's image")
	ErrPushFunctionImage  = errors.New("Couldn't push the function's image")
)

type Function struct {
	Path       string
	Version    string
	OldVersion string
	*Config
}

func newFunc() *Function {
	return &Function{}
}

type Config struct {
	Name        string `yaml:"name"`
	Private     bool   `yaml:"private"`
	ContentType string `yaml:"content_type"`
}

func newFuncConfig() *Config {
	return &Config{}
}

func getFunction(p string) (*Function, error) {
	cfg, err := getConfig(p)
	if err != nil {
		return nil, err
	}

	version := getVersion(p)
	if version == "" {
		return nil, ErrInvalidVersion
	}

	fn := newFunc()
	fn.Version = version
	fn.Config = cfg
	fn.Path = p

	return fn, nil
}

func getConfig(path string) (*Config, error) {
	funcFile := filepath.Join(path, functionsYML)

	data, err := ioutil.ReadFile(funcFile)
	if err != nil {
		return nil, ErrLoadFunctionConfig
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getVersion(path string) string {
	versionFile := filepath.Join(path, "VERSION")

	data, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return ""
	}

	return strings.Split(string(data), "\n")[0]
}

func (fn *Function) bump() error {
	vfile := filepath.Join(fn.Path, "VERSION")

	s, err := storage.NewVersionStorage("file", initialVersion)
	version := bump.NewSemverBumper(s, vfile)
	newver, err := version.BumpPatchVersion("", "")
	if err != nil {
		return err
	}

	fn.OldVersion = fn.Version
	fn.Version = newver.String()

	ioutil.WriteFile(vfile, []byte(fn.Version), 0666)
	return nil
}

func (fn *Function) build() error {
	err := exec.Command("docker", "build", "-t", fmt.Sprintf("%s:%s", fn.Name, fn.Version), fn.Path).Run()
	if err != nil {
		return ErrBuildFunctionImage
	}
	return nil
}

func (fn *Function) push() error {
	err := exec.Command("docker", "push", fmt.Sprintf("%s:%s", fn.Name, fn.Version)).Run()
	if err != nil {
		return ErrPushFunctionImage
	}
	return nil
}

var ErrAPILoadRoutes = errors.New("")

func (fn *Function) updateRoutesImage() error {
	u := url.URL{
		Scheme:   "http",
		Host:     endpoint,
		Path:     "/v1/routes",
		RawQuery: "images=" + fn.Name + ":" + fn.Version,
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return ErrAPILoadRoutes
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrAPILoadRoutes
	}

	var rwrapper *funcs.RoutesWrapper
	json.Unmarshal(respData, &rwrapper)

	for _, route := range rwrapper.Routes {
		upd := &funcs.RouteWrapper{
			Route: &funcs.Route{
				Image: fmt.Sprintf("%s:%s", fn.Name, fn.Version),
			},
		}

		data, err := json.Marshal(upd)
		if err != nil {
			fmt.Printf("failed to update route `%s%s`\n", route.AppName, route.Path)
			continue
		}

		u.Path = fmt.Sprintf("/v1/apps/%s/routes%s", route.AppName, route.Path)
		u.RawQuery = ""

		req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("failed to update route `%s%s`\n", route.AppName, route.Path)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("failed to update route `%s%s`\n", route.AppName, route.Path)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("updated `%s%s` to image `%s:%s`\n", route.AppName, route.Path, fn.Name, fn.Version)
		} else {
			fmt.Printf("failed to update route `%s%s`\n", route.AppName, route.Path)
		}
	}

	return nil
}
