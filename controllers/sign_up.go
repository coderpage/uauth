package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"uauth/models"
	"uauth/storage"
	"uauth/tools/mail"
	"uauth/tools/secure"

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

	user := &models.User{}
	rqsBd := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &rqsBd)
	if err != nil {
		log.Error("Parse Request Body Err:", err)
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user.Email = rqsBd["Email"]
	user.Password = rqsBd["Password"]
	// 邮件激活链接，验证后跳转链接
	activeRedirect := rqsBd["redirect"]
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

		sendActiveEmail(user, activeRedirect)
		return
	}

	// errCode == 3, user exist
	if err == storage.ErrRowExist {
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
// redirectUrl uAuth 验证后重定向的地址
func sendActiveEmail(user *models.User, redirectUrl string) (err error) {

	activeToken := secure.GenerateToken(32)
	expire := time.Now().Add(24 * time.Hour)
	auth := &models.Auth{Uid: user.Id, Token: activeToken, Type: models.AuthTypeUserActive, ExpiryDate: expire}

	_, err = storage.AddNewAuth(auth)
	if err != nil {
		return errors.New("Add Auth Err:" + err.Error())
	}

	mailer, err := mail.NewServiceMailer()
	if err != nil {
		log.Error("Send Active Email Err:", err)
		return err
	}

	activeUrl := beego.AppConfig.String("serverbaseurl") + "/uauth/user/active?active=" + activeToken + "&redirect=" + redirectUrl
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
