package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// User
	UserGroupNoActived = "nonactived" // 未激活状态
	UserGroupCommon    = "common"     // 一般
	UserGropAdmin      = "admin"      // 管理员
	// Auth
	AuthTypeUserActive   = "u-active"   // 激活类型 Token
	AuthTypeUserSignIn   = "u-signIn"   // 登录类型 Token
	AuthTypeUserFindPwd  = "u-findPwd"  // 找回密码类型 Token
	AuthTypeUserResetPwd = "u-resetPwd" // 重置密码类型 Token
)

func init() {
	orm.RegisterModel(new(User), new(Auth))
}

func Register() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	mysqlUser := beego.AppConfig.String("mysqluser")
	mysqlDb := beego.AppConfig.String("mysqldb")
	mysqlPwd := beego.AppConfig.String("mysqlpass")

	orm.RegisterDataBase("default", "mysql", mysqlUser+":"+mysqlPwd+"@/"+mysqlDb+"?charset=utf8&loc=Local")

	// 开启 ORM 调试模式
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)
}
