package functions

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list [target]",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		funcs := findFunctions(path)
		for _, fn := range funcs {
			cfg, _ := getConfig(fn)
			version := getVersion(fn)
			rel, _ := filepath.Rel(cwd, fn)
			fmt.Printf("* %s v%s (%s) \n", cfg.Name, version, rel)
		}
	},
}
