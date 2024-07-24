package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show crond command version info.",

	Run: func(cmd *cobra.Command, args []string) {
		RunVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func RunVersion() {
	fmt.Println("version 1.1.0")
}
