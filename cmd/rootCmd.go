package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "定时任务启动器",
	Long:  `定时任务启动器`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("RootCmd called")
		RunCrond()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	loggerPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	loggerPath = loggerPath + "/log/" + time.Now().Format("200601/")
	_, err := os.Stat(loggerPath) // os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			err := os.MkdirAll(loggerPath, 0755)
			if err != nil {
				fmt.Println("MkdirAll failed ")
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
