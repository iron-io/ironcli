package functions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	ErrInitFunctionExists = errors.New(functionsYML + " already exists on that function")
)

var InitCmd = &cobra.Command{
	Use: "init [target]",
	Run: func(cmd *cobra.Command, args []string) {
		path := ""
		if len(args) > 0 {
			if filepath.IsAbs(args[0]) {
				path = filepath.Clean(args[0])
			} else {
				path = filepath.Join(cwd, filepath.Clean(args[0]))
			}
		}
		err := initFn(path)
		if err == nil {
			fmt.Println("initialized", path)
		} else {
			fmt.Println("Error:", err)
		}
	},
}

func initFn(path string) error {

	funcFile := filepath.Join(path, functionsYML)
	versionFile := filepath.Join(path, "VERSION")

	_, err := os.Open(funcFile)
	if err == nil {
		return ErrInitFunctionExists
	}

	config := newFuncConfig()
	data, err := yaml.Marshal(&config)

	if err != nil {
		return err
	}

	if initName != "" {
		config.Name = initName
	}

	if initPrivate {
		config.Private = true
	}

	ioutil.WriteFile(funcFile, data, 0666)
	ioutil.WriteFile(versionFile, []byte(initialVersion), 0666)

	return nil
}
