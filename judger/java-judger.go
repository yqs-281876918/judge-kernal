package judger

import (
	"judge-kernel/util"
	"os/exec"
	"path/filepath"
	"strings"
)

type JavaJudger struct {
	JudgerBase
}

func (this *JavaJudger) JudgeCode(executablePath string) interface{} {
	//初始化cmd
	className := filepath.Base(executablePath)
	className = strings.Split(className, ".")[0]
	run_cmd := exec.Command("java", className)
	run_cmd.Dir = filepath.Dir(executablePath)
	return this.RunWithMemoryAndTimeoutMonitor(run_cmd)
}

func (this *JavaJudger) CompileCode(codePath string) (bool, string) {
	cmd := exec.Command("javac", "-encoding", "utf8", codePath)
	output_bytes, err := cmd.CombinedOutput()
	if err != nil {
		output := string(output_bytes)
		if util.IsGBK(output_bytes) {
			temp, err := util.GbkToUtf8(output_bytes)
			if err != nil {
				return false, "出现乱码，无法显示具体错误信息，请确保终端编码为UTF8或者GBK"
			}
			output = string(temp)
		}
		return false, output
	}
	fileName := filepath.Base(codePath)
	fileName = strings.Split(fileName, ".")[0]
	exePath := filepath.Dir(codePath) + "/" + fileName + ".class"
	return true, exePath
}
