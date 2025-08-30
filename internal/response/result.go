package response

import "log"

var (
	SUCCESS = 0
	FAIL    = 1000
)

type Result struct {
	code int
	msg  string
}

func (r *Result) Code() int {
	return r.code
}

func (r *Result) Msg() string {
	return r.msg
}

var (
	_code    = map[int]struct{}{}
	_message = make(map[int]Result)
)

func RegisterResult(code int, msg string) Result {
	_, ok := _code[code]
	if ok {
		log.Printf("业务码 %d 已存在", code)
		return _message[code]
	}

	result := Result{
		code: code,
		msg:  msg,
	}

	_code[code] = struct{}{}
	_message[code] = result

	return result
}

var (
	// 注册业务码
	Success = RegisterResult(SUCCESS, "成功")
	Fail    = RegisterResult(FAIL, "失败")

	ErrUserNotFound = RegisterResult(1000, "未知用户")
)
