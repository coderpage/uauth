package storage

import (
	"errors"

	"github.com/astaxie/beego/logs"
)

var (
	log         *logs.BeeLogger
	ErrRowExist = errors.New("this row is allready exist")
	ErrNoRows   = errors.New("no rows return")
)

func init() {
	log = logs.NewLogger(10000)
	// log.Async()
	log.SetLogger("console", "")
}
