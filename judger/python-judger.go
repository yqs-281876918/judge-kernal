package judger

import (
	"os/exec"
)

type PythonJudger struct {
	JudgerBase
}

func (this *PythonJudger) JudgeCode(executablePath string) int {
	run_cmd := exec.Command("python", executablePath)
	return this.RunWithMemoryAndTimeoutMonitor(run_cmd)
}
func (this *PythonJudger) CompileCode(codePath string) (bool, string) {
	return true, codePath
}
