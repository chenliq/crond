package cmd

import (
	"log"
	"os"

	. "fmt"
	"gopkg.in/yaml.v3"
)

var loggerPath string

type CfgData struct {
	App struct {
		Env string `yaml:"env"`
	}

	PHPCommand    string `yaml:"PHPCommand"`
	PHPScript     string `yaml:"PHPScript"`
	PythonCommand string `yaml:"PythonCommand"`
	PythonScript  string `yaml:"PythonScript"`

	Crond []CrondItem
}

type CrondItem struct {
	Name        string   `yaml:"name"`
	Args        []string `yaml:"args"`
	Type        string   `yaml:"type"`
	SpecTimer   string   `yaml:"specTimer"`
	Interpreter string   `yaml:"interpreter"`
	Script      string   `yaml:"script"`
}

var Config CfgData

var BasePath string

func init() {
	var err error

	parsePath()

	file := BasePath + "/config/crond.yaml"
	// Println(file)
	_, err = os.Stat(file)
	if err != nil {
		file = "./config/crond.yaml"
	}
	// Println(file)
	data, err := os.ReadFile(file)

	if err == nil {
		err = yaml.Unmarshal(data, &Config)
		// 当有异常时：
		if err != nil {
			Println("解析失败:", err)
			//	退出读取 code == 0时，表示读取成功  code == 1时，表示读取失败退出
			os.Exit(1)
		}
	} else {
		log.Printf("读取配置文件失败 #%v", err)
	}
}
