package functions

import "github.com/spf13/cobra"

var TestCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
