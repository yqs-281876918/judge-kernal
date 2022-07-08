package judge_err

const (
	Pass = iota
	UnknownError
	LostInputFile
	LostOutputFile
	WrongCode
	SystemException
	WrongAnswer
	Timeout
	OutOfMemory
)

type ErrorMsg struct {
	ErrorType string `json:"error_type"`
	Detail    string `json:"detail"`
}

func getDefaultMsg() *ErrorMsg {
	msg := new(ErrorMsg)
	msg.ErrorType = "unknown error"
	msg.Detail = "未知错误"
	return msg
}

func GetMsgByError(errorType int) *ErrorMsg {
	msg := getDefaultMsg()
	switch errorType {
	case Pass:
		msg.ErrorType = "pass"
		msg.Detail = "测试用例通过"
	case LostInputFile:
		msg.ErrorType = "lost input file"
		msg.Detail = "该题目输入文件丢失"
	case LostOutputFile:
		msg.ErrorType = "lost output file"
		msg.Detail = "该题目输出文件丢失"
	case WrongCode:
		msg.ErrorType = "wrong code"
		msg.Detail = "代码存在执行错误"
	case SystemException:
		msg.ErrorType = "system exception"
		msg.Detail = "评测系统出错,请联系管理员"
	case WrongAnswer:
		msg.ErrorType = "wrong answer"
		msg.Detail = "测试用例不通过"
	case Timeout:
		msg.ErrorType = "timeout"
		msg.Detail = "程序执行时间超时"
	case OutOfMemory:
		msg.ErrorType = "out of memory"
		msg.Detail = "内存超出限制"
	}
	return msg
}
