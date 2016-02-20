package main

import (
	"uauth/models"
	"uauth/routers"

	"github.com/astaxie/beego"
)

func init() {
	models.Register()
	routers.Register()
}

func main() {
	beego.Run()
}
