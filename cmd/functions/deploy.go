package functions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	ErrFunctionNotChanged = errors.New("Function hasn't been changed")
)

var DeployCmd = &cobra.Command{
	Use: "deploy [target]",
	Run: func(cmd *cobra.Command, args []string) {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		funcs := findFunctions(path)
		for _, f := range funcs {
			vfile := filepath.Join(f, "VERSION")

			versionInfo, err := os.Stat(vfile)
			if err != nil {
				fmt.Println("Coudn't find VERSION for", f)
				continue
			}

			changed := false
			filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
				if path != vfile && !info.IsDir() && !versionInfo.ModTime().After(info.ModTime()) {
					changed = true
					return nil
				}
				return nil
			})

			if !changed {
				fmt.Println("error", ErrFunctionNotChanged)
				continue
			}

			fn, err := getFunction(f)
			if err != nil {
				fmt.Println("error", f, ErrLoadFunctionConfig)
				continue
			}
			fmt.Println("found function", fn.Name, f)

			err = fn.bump()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("bumped", fn.Name, "version")

			err = fn.build()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("builded", fn.Name, "image")

			err = fn.push()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("pushed", fn.Name, "image")

			err = fn.updateRoutesImage()
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
		}
	},
}
