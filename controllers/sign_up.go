package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"uauth/controllers/mail"
	"uauth/models"
	"uauth/storage"

	"github.com/astaxie/beego"
)

type SignUpHandler struct {
	BaseController
}

func (this *SignUpHandler) SignUp() {
	resp := make(Response)

	body := this.Ctx.Request.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error("Read Request Body Err:", err)
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Read Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user := new(models.User)
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		log.Error("Parse Request Body Err:", err)
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if user.Email == "" || user.Password == "" {
		log.Error("Parse Request Body Err:", err)

		if user.Email == "" {
			resp.SetStatus(StatusInvalidUserName)
		} else {
			resp.SetStatus(StatusInvalidPwd)
		}
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	created, err := storage.CreateUser(user)
	if created {
		log.Info("Create New User Success:", user.String())
		resp.SetStatus(http.StatusOK)
		resp.SetMessage("OK")
		this.Data["json"] = resp
		this.ServeJSON()
		sendActiveEmail(user, "abc", "#")
		return
	}

	errMsg := strings.Split(err.Error(), "-")
	errCode := errMsg[0]
	// errCode == 3, user exist
	if errCode == "3" {
		log.Error("Create New User Err: User Exist")
		resp.SetStatus(StatusUserExist)
		resp.SetMessage("User Exist")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	log.Error("Create New User Err: ", err)
	resp.SetStatus(StatusUkownError)
	resp.SetMessage("Create User Failed")
	this.Data["json"] = resp
	this.ServeJSON()
}

// sendActiveEmail 发送激活邮件
// activeToken 激活码
// redirectUrl uAuth 验证后重定向的地址
func sendActiveEmail(user *models.User, activeToken, redirectUrl string) (err error) {
	mailer, err := mail.NewServiceMailer()
	if err != nil {
		log.Error("Send Active Email Err:", err)
		return err
	}

	activeUrl := beego.AppConfig.String("serverbaseurl") + "/user/active?active=" + activeToken + "&redirect=" + redirectUrl
	body := fmt.Sprintf(`	尊敬的 %s 您好！
<br>
点击 <a href="%s">链接</a> 可激活您的的账号！
<br>
为保障您的帐号安全，请在24小时内点击该链接，您也可以将链接复制到浏览器地址栏访问。如果您并未尝试激活邮箱，请忽略本邮件，由此给您带来的不便请谅解。
<br>
<br>
本邮件由系统自动发出，请勿直接回复！
<br>
<br>
`, user.Email, activeUrl)
	err = mailer.SendMail(user.Email, "uAuth", "请激活账号", "html", body)
	if err != nil {
		log.Error("Send Active Email Err:", err)
	}

	return err
}
