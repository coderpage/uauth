package controllers

const (
	StatusUkownError          = -1  // 未知错误
	StatusUnprocessableEntity = 422 // 无法处理的请求实体
	StatusInvalidUserName     = 430 // 无效用户名
	StatusInvalidPwd          = 431 // 无效密码
	StatusUserExist           = 432 // 该用户已存在
	StatusUserNotExist        = 433 // 用户不存在
	StatusUserNotActivated    = 434 // 用户未激活
	StatusWrongUserNameOrPwd  = 435 // 用户名或密码错误
)

var statusText = map[int]string{
	StatusUkownError:          "Unkown Error",
	StatusUnprocessableEntity: "Parse Request Entity Error",
	StatusInvalidUserName:     "Invalid User Name",
	StatusInvalidPwd:          "Invalid Password",
	StatusUserExist:           "User Already Exist",
	StatusUserNotExist:        "User Not Exist",
	StatusUserNotActivated:    "User Not Activated",
	StatusWrongUserNameOrPwd:  "User Name Or Password Is Wrong",
}

func StatusText(code int) string {
	return statusText[code]
}
