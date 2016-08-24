package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	RunZIP  string
	RunName string
)

var RunCmd = &cobra.Command{
	Use: "run [image]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}
