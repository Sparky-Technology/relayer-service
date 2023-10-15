package utils

/**
 * @Description //TODO 定义业务状态码
 **/

type MyCode int64

const (
	CodeSuccess               MyCode = 1000
	CodeInvalidParams         MyCode = 1001
	CodeNotPassFailGuardCheck MyCode = 1002
	CodeSendTxFail            MyCode = 1003

	CodeUserExist         MyCode = 9002
	CodeUserNotExist      MyCode = 9003
	CodeInvalidPassword   MyCode = 9004
	CodeServerBusy        MyCode = 9005
	CodeInvalidToken      MyCode = 9006
	CodeInvalidAuthFormat MyCode = 9007
	CodeNotLogin          MyCode = 9008
)

var msgFlags = map[MyCode]string{
	CodeSuccess:               "success",
	CodeInvalidParams:         "invalid params",
	CodeNotPassFailGuardCheck: "not pass FailGuard check",
	CodeSendTxFail:            "send tx fail",

	CodeUserExist:         "用户名重复",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "用户名或密码错误",
	CodeServerBusy:        "服务繁忙",
	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
