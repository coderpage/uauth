package storage

import (
	"github.com/astaxie/beego/logs"
)

var log *logs.BeeLogger

func init() {
	log = logs.NewLogger(10000)
	// log.Async()
	log.SetLogger("console", "")
}
