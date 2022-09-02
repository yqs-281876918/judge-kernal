package judger

import (
	"os/exec"
)

type GoJudger struct {
	JudgerBase
}

func (this *GoJudger) JudgeCode(executablePath string) interface{} {
	run_cmd := exec.Command("go", "run", executablePath)
	return this.RunWithMemoryAndTimeoutMonitor(run_cmd)
}
func (this *GoJudger) CompileCode(codePath string) (bool, string) {
	return true, codePath
}
