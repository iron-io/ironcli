package functions

import (
	"os"
	"path/filepath"

	"github.com/iron-io/ironcli/cmd/functions/apps"
	"github.com/iron-io/ironcli/cmd/functions/routes"
	"github.com/spf13/cobra"
)

var (
	cwd, _         = os.Getwd()
	endpoint       = ""
	findAllFlag    bool
	functionsYML   = "functions.yml"
	initialVersion = "0.0.1"

	initName    = ""
	initPrivate = true
)

var RootCmd = &cobra.Command{
	Use:     "fn",
	Aliases: []string{"functions"},
}

func findFunctions(path string) []string {
	wd := cwd
	if path != "" {
		wd = filepath.Join(cwd, path)
	}

	funcDirs := []string{}
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Valid function directory should have at least this two files.
			_, err := os.Stat(filepath.Join(path, "VERSION"))
			_, err2 := os.Stat(filepath.Join(path, functionsYML))

			if err == nil && err2 == nil {
				funcDirs = append(funcDirs, path)
			}

			if !findAllFlag {
				return filepath.SkipDir
			}
		}

		return nil
	})

	return funcDirs
}

func init() {
	RootCmd.AddCommand(
		apps.AppCmd,
		routes.RouteCmd,
		DeployCmd,
		ListCmd,
		InitCmd,
	)
	if os.Getenv("FUNCTIONS_ENDPOINT") != "" {
		endpoint = os.Getenv("FUNCTIONS_ENDPOINT")
	}
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "E", "127.0.0.1:8080", "functions API endpoint")
	RootCmd.PersistentFlags().BoolVarP(&findAllFlag, "all", "a", false, "select all functions")

	InitCmd.Flags().StringVarP(&initName, "name", "n", "", "define function name")
	InitCmd.Flags().BoolVarP(&initPrivate, "private", "p", false, "define if functions is private")

}
