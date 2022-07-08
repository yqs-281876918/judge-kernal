package judger

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type JavaJudger struct {
	JudgerBase
}

func (this *JavaJudger) JudgeCode(executablePath string) int {
	//初始化cmd
	className := filepath.Base(executablePath)
	className = strings.Split(className, ".")[0]
	run_cmd := exec.Command("java", className)
	run_cmd.Dir = filepath.Dir(executablePath)
	return this.RunWithMemoryAndTimeoutMonitor(run_cmd)
}

func (this *JavaJudger) CompileCode(codePath string) (bool, string) {
	cmd := exec.Command("javac", "-encoding", "utf8", codePath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	fileName := filepath.Base(codePath)
	fileName = strings.Split(fileName, ".")[0]
	exePath := filepath.Dir(codePath) + "/" + fileName + ".class"
	return true, exePath
}
