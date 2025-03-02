package ginx

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success() Result {
	return Result{
		Code: 0,
		Msg:  "请求成功",
	}
}

func InvalidToken() Result {
	return Result{
		Code: 2,
		Msg:  "无效的token",
	}
}

func InvalidParam() Result {
	return Result{
		Code: 4,
		Msg:  "参数错误",
	}
}

func InvalidParamWithMsg(msg string) Result {
	return Result{
		Code: 4,
		Msg:  msg,
	}
}

func SuccessWithData(data any) Result {
	return Result{
		Code: 0,
		Msg:  "请求成功",
		Data: data,
	}
}

func SuccessWithMsg(msg string) Result {
	return Result{
		Code: 0,
		Msg:  msg,
	}
}

func BizError(msg string) Result {
	return Result{
		Code: 1,
		Msg:  msg,
	}
}

func InternalError() Result {
	return Result{
		Code: 500,
		Msg:  "系统错误",
	}
}

func InternalErrorWithMsg(msg string) Result {
	return Result{
		Code: 500,
		Msg:  msg,
	}
}

func NotFound() Result {
	return Result{
		Code: 404,
		Msg:  "资源不存在",
	}
}

// i18n
var langCodeMsgMap = map[string]map[int]string{
	"zh": {
		0:   "请求成功",
		1:   "业务错误",
		2:   "无效的token",
		4:   "参数错误",
		500: "系统错误",
		404: "资源不存在",
	},
	"en": {
		0:   "Request Success",
		1:   "Biz Error",
		2:   "Invalid Token",
		4:   "Invalid Param",
		500: "Internal Error",
		404: "Resource Not Found",
	},
}
