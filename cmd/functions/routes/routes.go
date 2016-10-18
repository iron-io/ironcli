package routes

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
	ErrMissingAppName   = errors.New("Missing app's name")
	ErrMissingRouteName = errors.New("Missing route's name")
	ErrMissingImage     = errors.New("Missing route's image")
	ErrCreateRoute      = errors.New("Could not create the route on given endpoint")
	ErrListRoute        = errors.New("Could list routes on given endpoint")
	ErrDeleteRoute      = errors.New("Could delete route on given endpoint")
)

var RouteCmd = &cobra.Command{Use: "routes"}

var routeCreateCmd = &cobra.Command{
	Use:     "create APPNAME PATH IMAGE",
	Aliases: []string{"put"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrMissingAppName
		}

		if len(args) <= 1 {
			return ErrMissingRouteName
		}

		if len(args) <= 2 {
			return ErrMissingImage
		}

		appname := args[0]
		path := args[1]
		image := args[2]

		upd := &funcs.RouteWrapper{
			Route: &funcs.Route{
				Path:    path,
				AppName: appname,
				Image:   image,
			},
		}

		data, err := json.Marshal(upd)
		if err != nil {
			return ErrCreateRoute
		}

		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps/" + appname + "/routes",
		}

		req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(data))
		if err != nil {
			return ErrCreateRoute
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrCreateRoute
		}

		if resp.StatusCode != http.StatusOK {
			return ErrCreateRoute
		}

		fmt.Printf("created route `%s%s` on `%s`\n", appname, path, endpoint)
		return nil
	},
}

var routeListCmd = &cobra.Command{
	Use:     "list APPNAME",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrMissingAppName
		}

		appname := args[0]

		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps/" + appname + "/routes",
		}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return ErrListRoute
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrListRoute
		}

		if resp.StatusCode != http.StatusOK {
			return ErrListRoute
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ErrListRoute
		}

		var wroutes *funcs.RoutesWrapper
		json.Unmarshal(data, &wroutes)
		if err != nil {
			return ErrListRoute
		}

		for _, route := range wroutes.Routes {
			fmt.Printf("* %s on %s (http://%s/r/%s%s)\n", route.Path, appname, endpoint, appname, route.Path)
		}

		return nil
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:     "delete APPNAME PATH",
	Aliases: []string{"del"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrMissingAppName
		}

		if len(args) <= 1 {
			return ErrMissingRouteName
		}

		appname := args[0]
		path := args[1]

		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		u := &url.URL{
			Scheme: "http",
			Host:   endpoint,
			Path:   "/v1/apps/" + appname + "/routes" + path,
		}

		req, err := http.NewRequest("DELETE", u.String(), nil)
		if err != nil {
			return ErrDeleteRoute
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ErrDeleteRoute
		}

		if resp.StatusCode != http.StatusOK {
			return ErrDeleteRoute
		}

		fmt.Printf("deleted route `%s%s` from `%s`\n", appname, path, endpoint)
		return nil
	},
}

func init() {
	RouteCmd.AddCommand(
		routeCreateCmd,
		routeListCmd,
		routeDeleteCmd,
	)
}
