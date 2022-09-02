package judge_err

import "judge-kernel/global"

type ResultMsg struct {
	MsgType     string `json:"msg_type"`
	Description string `json:"description"`
}

type PassMsg struct {
	ResultMsg
	TimeCost   int     `json:"time_cost"`
	MemoryCost float64 `json:"memory_cost"`
}

func CreatePassMsg(time int, memory float64) *PassMsg {
	msg := &PassMsg{}
	msg.MsgType = "pass"
	msg.Description = "测试用例通过"
	msg.TimeCost = time
	msg.MemoryCost = memory
	return msg
}

type RuntimeErrorMsg struct {
	ResultMsg
	Detail string `json:"detail"`
}

func CreateRuntimeErrorMsg(errorDetail string) *RuntimeErrorMsg {
	msg := &RuntimeErrorMsg{}
	msg.MsgType = "runtime-error"
	msg.Description = "代码存在运行时异常"
	msg.Detail = errorDetail
	return msg
}

type SystemExceptionMsg struct {
	ResultMsg
	Detail string `json:"detail"`
}

func CreateSystemExceptionMsg(detail string) *SystemExceptionMsg {
	msg := &SystemExceptionMsg{}
	msg.MsgType = "system-exception"
	msg.Description = "评测系统出现错误,请稍后重试"
	msg.Detail = detail
	return msg
}

type WrongAnswerMsg struct {
	ResultMsg
	Input        string `json:"input"`
	ExpectOutput string `json:"expect_output"`
	ActualOutput string `json:"actual_output"`
}

func CreateWrongAnswerMsg(input string, expect string, actual string) *WrongAnswerMsg {
	msg := &WrongAnswerMsg{}
	msg.MsgType = "wrong-answer"
	msg.Description = "测试用例不通过"
	msg.Input = input
	msg.ExpectOutput = expect
	msg.ActualOutput = actual
	return msg
}

type TimeoutMsg struct {
	ResultMsg
	LimitedTime int `json:"limited_time"`
}

func CreateTimeoutMsg() *TimeoutMsg {
	msg := &TimeoutMsg{}
	msg.MsgType = "timeout"
	msg.Description = "运行时间超限"
	msg.LimitedTime = global.Arguments.TimeOut
	return msg
}

type OutOfMemoryMsg struct {
	ResultMsg
	LimitedMemorySize int `json:"limited_memory_size"`
}

func CreateOutOfMemoryMsg() *OutOfMemoryMsg {
	msg := &OutOfMemoryMsg{}
	msg.MsgType = "out-of-memory"
	msg.Description = "运行内存超限"
	msg.LimitedMemorySize = global.Arguments.MemoryLimit
	return msg
}

type CompileFailedMsg struct {
	ResultMsg
	Detail string `json:"detail"`
}

func CreateCompileFailedMsg(detail string) *CompileFailedMsg {
	msg := &CompileFailedMsg{}
	msg.MsgType = "compile-failed"
	msg.Description = "编译错误"
	msg.Detail = detail
	return msg
}
