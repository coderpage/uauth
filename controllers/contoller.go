package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type Response map[string]interface{}

func (resp Response) JsonString() string {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Error("Marshal Response Err:", err)
	}
	return string(respBytes)
}

func (resp Response) SetStatus(status interface{}) {
	resp["Status"] = status
}

func (resp Response) SetMessage(message interface{}) {
	resp["Message"] = message
}

const (
	StatusUnprocessableEntity = 422 // 无法处理的请求实体
	StatusInvalidUserName     = 430 // 无效用户名
	StatusInvalidPwd          = 431 // 无效密码
	StatusUserExist           = 432 // 该用户已存在
	StatusUserNotExist        = 433 // 用户不存在
	StatusUkownError          = 434 // 未知错误
)

var log *logs.BeeLogger

func init() {
	log = logs.NewLogger(10000)
	log.Async()
	log.SetLogger("console", "")
	log.EnableFuncCallDepth(true)
}

type BaseController struct {
	beego.Controller
}
