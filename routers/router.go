package routers

import (
	"uauth/controllers"

	"github.com/astaxie/beego"
)

func Register() {
	beego.Router("/uauth/signup", &controllers.SignUpHandler{}, "post:SignUp")
	beego.Router("/uauth/signin", &controllers.SignInHandler{}, "post:SignIn")
}
