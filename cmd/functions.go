package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "fmt"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

func NewLogger(args ...string) *log.Logger {
	var fileName string
	if len(args) > 0 {
		fileName = args[0]
	} else {
		fileName = time.Now().Format("2")
	}
	file := loggerPath + fileName + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.Ldate|log.Ltime)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	return logger
}

func PrintWithColor(data any, args ...string) {
	color := ""
	var format = "%s"
	if len(args) > 0 {
		color = args[0]
	}
	if len(args) > 1 {
		format = args[1]
	}

	switch color {
	case "red":
		Printf("\033[1;31;40m"+format+"\033[0m\n", data)
	case "green":
		Printf("\033[1;32;40m"+format+"\033[0m\n", data)
	case "yellow":
		Printf("\033[1;33;40m"+format+"\033[0m\n", data)
	case "blue":
		Printf("\033[1;34;40m"+format+"\033[0m\n", data)
	case "grey":
		Printf("\033[1;30;40m"+format+"\033[0m\n", data)
	default:
		Printf("\033[1;37;40m"+format+"\033[0m\n", data)
	}
}

func parsePath() {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			BasePath = path.Dir(path.Dir(filename))
		} else {
			panic("get basePath file")
		}
	} else {
		BasePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		// Println(basePath)
	}
}

func parseFilePath(file string) string {
	if strings.HasPrefix(file, "/") {
		return file
	} else if strings.HasPrefix(file, "./") {
		return BasePath + file[1:]
	} else if strings.HasPrefix(file, "../") {
		return path.Dir(BasePath) + file[2:]
	} else {
		return BasePath + "/" + file
	}
}
