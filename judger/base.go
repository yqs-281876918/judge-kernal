package judger

import (
	"bytes"
	"github.com/shirou/gopsutil/process"
	judge_err "judge-kernel/err"
	"judge-kernel/global"
	"judge-kernel/util"
	"os"
	"os/exec"
	"time"
)

type Judger interface {
	JudgeCode(string) int
	CompileCode(string) (bool, string)
}

type JudgerBase struct {
}

func (this *JudgerBase) ReadInputBytes() ([]byte, error) {
	content, err := os.ReadFile(global.Arguments.InputPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (this *JudgerBase) readOutputBytes() ([]byte, error) {
	content, err := os.ReadFile(global.Arguments.OutputPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (this *JudgerBase) RunWithMemoryAndTimeoutMonitor(run_cmd *exec.Cmd) int {
	outFlagChan := make(chan struct{}, 1)
	errorChan := make(chan int, 1)
	go this.run(run_cmd, errorChan)
	go this.startMemoryMonitor(run_cmd, outFlagChan)
	select {
	case <-outFlagChan:
		return judge_err.OutOfMemory
	case <-time.After(time.Duration(global.Arguments.TimeOut * 1e6)):
		return judge_err.Timeout
	case code := <-errorChan:
		return code
	}
}

func (this *JudgerBase) startMemoryMonitor(run_cmd *exec.Cmd, errorChan chan<- struct{}) {
	for run_cmd.Process == nil {
	}
	p, err := process.NewProcess(int32(run_cmd.Process.Pid))
	if err != nil {
		return
	}
	for {
		memInfo, err := p.MemoryInfo()
		if err != nil {
			return
		}
		MB := float32(memInfo.RSS) / (1024 * 1024)
		if MB >= float32(global.Arguments.MemoryLimit) {
			errorChan <- struct{}{}
			return
		}
		time.Sleep(time.Millisecond * 1)
	}
}

func (this *JudgerBase) run(run_cmd *exec.Cmd, errorChan chan int) {
	//获取输入
	in_content, err := this.ReadInputBytes()
	if err != nil {
		errorChan <- judge_err.LostInputFile
		return
	}
	//去掉空行结尾
	in_content = bytes.TrimRight(in_content, "\n")
	in_content = bytes.TrimRight(in_content, "\r")
	//写入输入
	in_writer, err := run_cmd.StdinPipe()
	if err != nil {
		errorChan <- judge_err.SystemException
		return
	}
	_, err = in_writer.Write(in_content)
	if err != nil {
		errorChan <- judge_err.SystemException
		return
	}
	err = in_writer.Close()
	if err != nil {
		errorChan <- judge_err.SystemException
		return
	}
	//获取程序输出
	outputBytes, err := run_cmd.CombinedOutput()
	if err != nil {
		errorChan <- judge_err.WrongCode
		return
	}
	outputBytes = bytes.TrimRight(outputBytes, "\n")
	outputBytes = bytes.TrimRight(outputBytes, "\r")
	var possibleCharset = "utf8"
	if util.IsGBK(outputBytes) {
		possibleCharset = "gbk"
	}
	exec_output := util.ConvertToString(string(outputBytes), possibleCharset, "utf8")
	//获取目标输出
	target_output_bytes, err := this.readOutputBytes()
	if err != nil {
		errorChan <- judge_err.LostOutputFile
		return
	}
	target_output_bytes = bytes.TrimRight(target_output_bytes, "\n")
	target_output_bytes = bytes.TrimRight(target_output_bytes, "\r")
	target_output := string(target_output_bytes)
	if exec_output == target_output {
		errorChan <- judge_err.Pass
		return
	} else {
		errorChan <- judge_err.WrongAnswer
		return
	}
}
