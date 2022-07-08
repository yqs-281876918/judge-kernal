package judger

import (
	"judge-kernel/global"
	"os/exec"
	"path/filepath"
	"strings"
)

type JavaJudger struct {
	JudgerBase
}

func (this *JavaJudger) JudgeCode() int {
	//初始化cmd
	className := filepath.Base(global.Arguments.ExecutablePath)
	className = strings.Split(className, ".")[0]
	run_cmd := exec.Command("java", className)
	run_cmd.Dir = filepath.Dir(global.Arguments.ExecutablePath)
	return this.RunWithMemoryAndTimeoutMonitor(run_cmd)
}
