package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	judge_err "judge-kernel/err"
	"judge-kernel/global"
	jd "judge-kernel/judger"
	"os"
	"path/filepath"
	"strings"
)

func initArg() {
	var language string
	var input_path string
	var output_path string
	var result_path string
	var executable_path string
	var timeout int
	var memory_limit int
	flag.StringVar(&language, "lang", "java", "输入程序语言类型(JAVA/GO)")
	flag.StringVar(&input_path, "input", "/default.in", "输入输入文件路径")
	flag.StringVar(&output_path, "output", "/default.out", "输入输出文件路径")
	flag.StringVar(&result_path, "result", "/result.json", "输入结果文件路径")
	flag.StringVar(&executable_path, "executable", "/null", "输入可执行文件路径")
	flag.IntVar(&timeout, "timeout", 1000, "输入程序时间限制(ms)")
	flag.IntVar(&memory_limit, "memory-limit", 256, "输入程序内存限制(MB)")
	flag.Parse()
	global.Arguments.Language = strings.ToLower(language)
	global.Arguments.InputPath = input_path
	global.Arguments.OutputPath = output_path
	global.Arguments.ResultPath = result_path
	global.Arguments.ExecutablePath = executable_path
	global.Arguments.TimeOut = timeout
	global.Arguments.MemoryLimit = memory_limit
}

func initConfig() error {
	return ini.MapTo(global.AppConfig, "./app.ini")
}

func main() {
	initArg()
	err := initConfig()
	if err != nil {
		fmt.Println(err)
	}
	var judger jd.Judger
	switch global.Arguments.Language {
	case "java":
		judger = new(jd.JavaJudger)
	case "go":
		judger = new(jd.GoJudger)
	default:
		judger = new(jd.JavaJudger)
	}
	errorID := judger.JudgeCode()
	errMsg := judge_err.GetMsgByError(errorID)
	bytes, _ := json.Marshal(errMsg)
	_, err = os.Stat(global.Arguments.ResultPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(global.Arguments.ResultPath), os.ModePerm)
		if err != nil {
			fmt.Printf("%v", err)
		}
		_, err = os.Create(filepath.Base(global.Arguments.ResultPath))
	}
	err = ioutil.WriteFile(global.Arguments.ResultPath, bytes, os.ModePerm)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
