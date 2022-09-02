package judger

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/process"
	judge_err "judge-kernel/err"
	"judge-kernel/global"
	"judge-kernel/util"
	"math"
	"os"
	"os/exec"
	"time"
)

type Judger interface {
	JudgeCode(string) interface{}
	CompileCode(string) (bool, string)
}

type JudgerBase struct {
	maxMemory float64
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

func (this *JudgerBase) RunWithMemoryAndTimeoutMonitor(run_cmd *exec.Cmd) interface{} {
	outFlagChan := make(chan struct{}, 1)
	errorChan := make(chan interface{}, 1)
	go this.run(run_cmd, errorChan)
	go this.startMemoryMonitor(run_cmd, outFlagChan)
	select {
	case <-outFlagChan:
		if run_cmd.Process != nil {
			_ = run_cmd.Process.Kill()
		}
		return judge_err.CreateOutOfMemoryMsg()
	case <-time.After(time.Duration(global.Arguments.TimeOut * 1e6)):
		if run_cmd.Process != nil {
			_ = run_cmd.Process.Kill()
		}
		return judge_err.CreateTimeoutMsg()
	case code := <-errorChan:
		if run_cmd.Process != nil {
			_ = run_cmd.Process.Kill()
		}
		return code
	}
}

func (this *JudgerBase) startMemoryMonitor(run_cmd *exec.Cmd, errorChan chan<- struct{}) {
	this.maxMemory = 0
	for run_cmd.Process == nil {
		time.Sleep(time.Millisecond * 1)
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
		MB := float64(memInfo.RSS) / (1024 * 1024)
		if MB >= float64(global.Arguments.MemoryLimit) {
			errorChan <- struct{}{}
			return
		}
		this.maxMemory = math.Max(this.maxMemory, MB)
		time.Sleep(time.Millisecond * 5)
	}
}

func (this *JudgerBase) run(run_cmd *exec.Cmd, errorChan chan interface{}) {
	//获取输入
	in_content, err := this.ReadInputBytes()
	if err != nil {
		errorChan <- judge_err.CreateSystemExceptionMsg(fmt.Sprintf("无法读取输入文件 %s", global.Arguments.InputPath))
		return
	}
	//去掉空行结尾
	in_content = bytes.TrimRight(in_content, "\n")
	in_content = bytes.TrimRight(in_content, "\r")
	//获取目标输出
	target_output_bytes, err := this.readOutputBytes()
	if err != nil {
		errorChan <- judge_err.CreateSystemExceptionMsg(fmt.Sprintf("无法读取输出文件 %s", global.Arguments.OutputPath))
		return
	}
	target_output_bytes = bytes.TrimRight(target_output_bytes, "\n")
	target_output_bytes = bytes.TrimRight(target_output_bytes, "\r")
	target_output := string(target_output_bytes)
	//写入输入
	in_writer, err := run_cmd.StdinPipe()
	if err != nil {
		errorChan <- judge_err.CreateSystemExceptionMsg("打开输入流失败")
		return
	}

	//开始计时
	start := time.Now().UnixNano() / 1e6

	_, err = in_writer.Write(in_content)
	if err != nil {
		errorChan <- judge_err.CreateSystemExceptionMsg("写入输入失败")
		return
	}
	err = in_writer.Close()
	if err != nil {
		errorChan <- judge_err.CreateSystemExceptionMsg("关闭输入流失败")
		return
	}
	//获取程序输出
	outputBytes, err := run_cmd.CombinedOutput()
	if err != nil {
		if util.IsGBK(outputBytes) {
			temp, err := util.GbkToUtf8(outputBytes)
			if err != nil {
				errorChan <- judge_err.CreateSystemExceptionMsg("请确保系统终端编码为GBK/UTF-8")
			}
			errorChan <- judge_err.CreateRuntimeErrorMsg(string(temp))
		} else {
			errorChan <- judge_err.CreateRuntimeErrorMsg(string(outputBytes))
		}
		return
	}
	end := time.Now().UnixNano() / 1e6
	outputBytes = bytes.TrimRight(outputBytes, "\n")
	outputBytes = bytes.TrimRight(outputBytes, "\r")
	exec_output := string(outputBytes)
	if util.IsGBK(outputBytes) {
		temp, err := util.GbkToUtf8(outputBytes)
		if err != nil {
			errorChan <- judge_err.CreateSystemExceptionMsg("请确保系统终端编码为GBK/UTF-8")
		}
		exec_output = string(temp)
	}
	if exec_output == target_output {
		errorChan <- judge_err.CreatePassMsg(int(end-start), this.maxMemory)
		return
	} else {
		errorChan <- judge_err.CreateWrongAnswerMsg(string(in_content), target_output, exec_output)
		return
	}
}
