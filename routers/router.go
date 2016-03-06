package routers

import (
	"uauth/controllers"

	"github.com/astaxie/beego"
)

func Register() {
	beego.Router("/uauth/signup", &controllers.SignUpHandler{}, "post:SignUp")
	beego.Router("/uauth/signin", &controllers.SignInHandler{}, "post:SignIn")
	beego.Router("/uauth/user/active", &controllers.UserActiveHandler{}, "get:ActiveFromEmail")
	beego.Router("/uauth/user/active/sendemail", &controllers.UserActiveHandler{}, "post:ResendActivateEmail")
	beego.Router("/uauth/find/user/withtk", &controllers.UserDataHandler{}, "post:FindUserWithAuthToken")
	beego.Router("/uauth/user/fpwd/email", &controllers.ResetPwdHandler{}, "post:FindPwdByEmail")
	beego.Router("/uauth/user/fpwd/email", &controllers.ResetPwdHandler{}, "get:AuthResetAction")
	beego.Router("/uauth/user/resetpwd", &controllers.ResetPwdHandler{}, "post:ResetPwd")
}
