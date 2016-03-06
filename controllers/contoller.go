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

func (resp Response) SetData(name string, data interface{}) {
	resp[name] = data
}

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
