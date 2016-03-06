package models

import (
	"encoding/json"
	"time"
)

// User 用户表
type User struct {
	Id          int64     `json:"Id"`                              // 主键id
	UserName    string    `json:"UserName" orm:"size(100)"`        // 用户名
	Password    string    `json:"-"`                               // 用户密码
	Email       string    `json:"Email"`                           // 用户邮箱
	DisplayName string    `json:"DisplayName"`                     // 用户显示的名称
	Url         string    `json:"Url"`                             // 用户主页
	Created     time.Time `json:"Created" orm:"auto_now_add"`      // 用户注册的时间
	Activated   time.Time `json:"Activated"`                       // 用户最后活动的时间
	Logged      time.Time `json:"Logged"`                          // 用户上次登录的时间
	Group       string    `json:"Group" orm:"default(nonactived)"` // 用户组 nonactived | common | admin
}

func (this *User) String() string {
	// str := fmt.Sprintf("User: {Id:%d  UserName:%s  Password:%s  Email:%s  DisplayName:%s  Group:%s}", this.Id, this.UserName, this.Password, this.Email, this.DisplayName, this.Group)
	return this.JsonString()
}

func (this *User) JsonString() string {
	jbytes, err := json.Marshal(this)
	if err != nil {
		return "{}"
	}
	return string(jbytes)
}
