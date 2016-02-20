package mtest

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SER_PORT = "9090"
)

func RegisterMySQL() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	mysqlUser := "root"
	mysqlDb := "uauth"
	mysqlPwd := "62624"

	orm.RegisterDataBase("default", "mysql", mysqlUser+":"+mysqlPwd+"@/"+mysqlDb+"?charset=utf8")

	// 开启 ORM 调试模式
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)
}
