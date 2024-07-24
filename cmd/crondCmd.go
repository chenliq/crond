package cmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	. "fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(crondCmd)
}

var crondCmd = &cobra.Command{
	Use:   "crond",
	Short: "定时任务管理器.",

	// 运行命令
	Run: func(cmd *cobra.Command, args []string) {
		RunCrond()
	},
}

func RunCrond() {
	// 初始化配置信息
	// if Config.App.Env == "dev" {
	// config, _ := json.Marshal(Config)
	// PrintWithColor(string(config), "blue", "\nconfig: %s\n")
	// }
	PHPCommand := Config.PHPCommand
	PHPScript := Config.PHPScript
	PythonCommand := Config.PythonCommand
	PythonScript := Config.PythonScript
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	var data [][]string
	for _, value := range Config.Crond {

		// 命令行如果有需要更多参数，在配置文件中的 args 配置多条参数数组
		// args:
		//   - version
		//   - -l 5
		crondItem := value

		var err error
		var scriptArgs string
		if len(value.Args) > 0 && value.Name == value.Args[0] {
			scriptArgs = strings.Join(value.Args[1:], " ")
		} else {
			scriptArgs = strings.Join(value.Args, " ")
		}
		if crondItem.Type == "PHPCmd" || crondItem.Type == "PythonCmd" {
			var command string
			var script string
			if crondItem.Type == "PHPCmd" {
				command = PHPCommand
			} else {
				command = PythonCommand
			}
			if crondItem.Script != "" {
				script = crondItem.Script
			} else if crondItem.Type == "PHPCmd" {
				script = PHPScript
			} else {
				script = PythonScript
			}
			script = parseFilePath(script)
			PrintWithColor(script, "green")

			args := append([]string{script}, crondItem.Args...)

			data = append(data, []string{value.Name, value.SpecTimer, command, script, scriptArgs})

			_, err = c.AddFunc(crondItem.SpecTimer, func() {
				go definedCmd(command, crondItem, args...)
			})
		} else if crondItem.Type == "elseCmd" {
			if crondItem.Interpreter == "" || crondItem.Script == "" {
				continue
			}
			var script string
			script = parseFilePath(crondItem.Script)
			PrintWithColor(script, "green")

			command := crondItem.Interpreter
			scriptArgs = strings.Join(value.Args, " ")

			args := append([]string{script}, crondItem.Args...)

			data = append(data, []string{value.Name, value.SpecTimer, command, script, scriptArgs})

			_, err = c.AddFunc(crondItem.SpecTimer, func() {
				go elseCmd(command, crondItem, args...)
			})
		} else {
			println("Unknown crond item type")
			continue
		}
		if err != nil {
			Printf("%+q Add task failed. \n", crondItem)
			continue
		}

	}

	// 用表格形式输出所有计划任务
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Crond", "Interpreter", "Script", "Options"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.BgBlackColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgHiBlackColor},
	)

	table.AppendBulk(data)
	table.Render()

	c.Run()

}

func definedCmd(command string, crondItem CrondItem, args ...string) {
	// PrintWithColor(crondItem.Type, "green")

	execCommand(command, crondItem, args)
}

func elseCmd(command string, crondItem CrondItem, args ...string) {
	// PrintWithColor(crondItem.Type, "green")

	execCommand(command, crondItem, args)
}

func execCommand(command string, crondItem CrondItem, args []string) {
	// 开始时间
	Printf("[%s] start: %s %+q, task output：\n", Now(), command, args)
	cmd := exec.Command(command, args...)
	// output, err := cmd.CombinedOutput()
	// Printf("output %+v\n", output)
	// if err != nil {
	// 	Printf("error stdout %+v\n", err)
	// }
	stdout, e := cmd.StdoutPipe()
	stderr, ee := cmd.StderrPipe()
	if e != nil {
		Printf("error stdout %+v\n", e)
	}

	defer func(cmd *exec.Cmd) {
		err := cmd.Wait()
		if err != nil {
			PrintWithColor(err, "red", "cmd.Wait failed: %v\n")
			return
		}
	}(cmd)
	err := cmd.Start()
	if err != nil {
		Printf("cmd.Run(%v) failed: %v\n", crondItem.Args, err)
		return
	}

	if e == nil {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			PrintWithColor(in.Text(), "blue", "%+v")
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}
	if ee == nil {
		fileName := strings.Replace(crondItem.Name, ":", "-", -1)
		logger := NewLogger(fileName)
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			logger.Printf(" %+v", in.Text())
			PrintWithColor(in.Text(), "blue", "%+v")
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}

	// 完成时间
	Printf("[%s] end: %s %+q\n\n", Now(), command, args)
}
