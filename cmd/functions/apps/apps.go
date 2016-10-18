package apps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	funcs "github.com/iron-io/functions/api/models"
	"github.com/spf13/cobra"
)

var (
	ErrMissingAppName = errors.New("Missing app's name")
	ErrCreateApp      = errors.New("Could not create the app on given endpoint")
	ErrListApp        = errors.New("Could list apps on given endpoint")
	ErrDeleteApp      = errors.New("Could delete app on given endpoint")
)

var AppCmd = &cobra.Command{Use: "apps"}

var appCreateCmd = &cobra.Command{
	Use:     "create [app name]",
	Aliases: []string{"put"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrMissingAppName
		}

		appname := args[0]

		upd := &funcs.AppWrapper{
			App: &funcs.App{
				Name: appname,
			},
		}

		data, err := json.Marshal(upd)
		if err != nil {
			return ErrCreateApp
		}

		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps",
		}

		req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(data))
		if err != nil {
			return ErrCreateApp
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrCreateApp
		}

		if resp.StatusCode != http.StatusOK {
			return ErrCreateApp
		}

		fmt.Printf("created app `%s` on `%s`\n", appname, endpoint)
		return nil
	},
}

var appListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps",
		}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return ErrListApp
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrListApp
		}

		if resp.StatusCode != http.StatusOK {
			return ErrListApp
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ErrListApp
		}

		var wapps *funcs.AppsWrapper
		json.Unmarshal(data, &wapps)
		if err != nil {
			return ErrListApp
		}

		for _, app := range wapps.Apps {
			fmt.Printf("* %s\n", app.Name)
		}

		return nil
	},
}

var appDeleteCmd = &cobra.Command{
	Use:     "delete [app name]",
	Aliases: []string{"del"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrMissingAppName
		}

		appname := args[0]
		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps/" + appname,
		}

		req, err := http.NewRequest("DELETE", u.String(), nil)
		if err != nil {
			return ErrDeleteApp
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrDeleteApp
		}

		if resp.StatusCode != http.StatusOK {
			return ErrDeleteApp
		}

		fmt.Printf("deleted app `%s` from `%s`\n", appname, endpoint)
		return nil
	},
}

func init() {
	AppCmd.AddCommand(
		appCreateCmd,
		appListCmd,
		appDeleteCmd,
	)
}
