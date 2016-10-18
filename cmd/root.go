package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "iron",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println("IronCLI v" + Version)
			os.Exit(0)
		}
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var (
	EnvFlag       string
	ProjectIDFlag string
	TokenFlag     string
	versionFlag   bool
)

const (
	Version = "0.2.0"
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&EnvFlag, "env", "e", "", "provide specific dev environment")
	RootCmd.PersistentFlags().StringVarP(&ProjectIDFlag, "project-id", "p", "", "provide projectid")
	RootCmd.PersistentFlags().StringVarP(&TokenFlag, "token", "t", "", "provide projectid")
	RootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "provide projectid")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
