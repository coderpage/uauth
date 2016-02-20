package routers

import (
	"uauth/controllers"

	"github.com/astaxie/beego"
)

func Register() {
	beego.Router("/signup", &controllers.SignUpHandler{}, "post:SignUp")
}
