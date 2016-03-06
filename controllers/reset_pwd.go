package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"time"
	"uauth/models"
	"uauth/storage"
	"uauth/tools/mail"
	"uauth/tools/secure"
)

type ResetPwdHandler struct {
	BaseController
}

// FindPwdByEmail 向指定邮箱发送找回密码邮件
func (this *ResetPwdHandler) FindPwdByEmail() {
	resp := make(Response)

	body := this.Ctx.Request.Body
	defer body.Close()

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Read Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	rqsBd := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &rqsBd)
	if err != nil {
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	email := rqsBd["Email"]
	redirect := rqsBd["Redirect"]
	if email == "" || redirect == "" {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("miss attr Email or Redirect in body")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user, err := storage.FindUserByEmail(email)
	if err != nil {
		resp.SetStatus(StatusUserNotExist)
		resp.SetMessage("user not exist")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	err = sendFindPwdEmail(user, redirect)
	if err != nil {
		resp.SetStatus(StatusUkownError)
		resp.SetMessage("send email failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp.SetStatus(http.StatusOK)
	resp.SetMessage("Ok")
	this.Data["json"] = resp
	this.ServeJSON()
}

// AuthResetAction 响应找回密码邮件重置链接，验证此链接合法性，跳转请求服务器
func (this *ResetPwdHandler) AuthResetAction() {
	fpwdToken := this.GetString("reset")
	redirect := this.GetString("redirect")

	if fpwdToken == "" {
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusNotFound, "token not found")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("404! token not found!")
		}
		return
	}

	auth, err := storage.FindAuthByToken(fpwdToken)
	if err != nil {
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusNotFound, "token not found")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("404! token not found!")
		}
		return
	}

	if time.Now().After(auth.ExpiryDate) {
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusForbidden, "token is expired")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("403! token is expired!")
		}
		return
	}

	resetToken := secure.GenerateToken(32)
	expire := time.Now().Add(24 * time.Hour)
	auth = &models.Auth{Uid: auth.Uid, Token: resetToken, Server: redirect, Status: "ok", Type: models.AuthTypeUserResetPwd, ExpiryDate: expire}

	_, err = storage.AddNewAuth(auth)
	if err != nil {
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusInternalServerError, "generate reset password token failed")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("500! active failed!")
		}
		return
	}

	if redirect != "" {
		redirect = fmt.Sprintf(redirect+"?status=%d&reset=%s", http.StatusOK, resetToken)
		this.Redirect(redirect, 302)
	} else {
		this.Ctx.WriteString("200!")
	}
}

// ResetPwd 重置密码
func (this *ResetPwdHandler) ResetPwd() {
	resp := make(Response)

	body := this.Ctx.Request.Body
	defer body.Close()

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Read Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	rqsBd := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &rqsBd)
	if err != nil {
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resetToken := rqsBd["ResetToken"]
	newPwd := rqsBd["NewPassword"]
	if resetToken == "" || newPwd == "" {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("miss attr ResetToken or NewPassword in body")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	auth, err := storage.FindAuthByToken(resetToken)
	if err != nil {
		resp.SetStatus(http.StatusNotFound)
		resp.SetMessage("ResetToken not found")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if time.Now().After(auth.ExpiryDate) {
		resp.SetStatus(http.StatusForbidden)
		resp.SetMessage("ResetToken was expired")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if auth.Type != models.AuthTypeUserResetPwd {
		resp.SetStatus(http.StatusForbidden)
		resp.SetMessage("ResetToken miss type")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	newPwd, err = secure.DoMd5(newPwd)
	if err != nil {
		resp.SetStatus(http.StatusInternalServerError)
		resp.SetMessage("internal server error")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	user := &models.User{Id: auth.Uid, Password: newPwd}
	err = storage.UpdateUser(user, "Password")
	if err != nil {
		resp.SetStatus(http.StatusInternalServerError)
		resp.SetMessage("internal server error")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp.SetStatus(http.StatusOK)
	resp.SetMessage("reset password success")
	this.Data["json"] = resp
	this.ServeJSON()
}

// sendFindPwdEmail 发送重置密码认证邮件到邮箱
func sendFindPwdEmail(user *models.User, redirect string) (err error) {
	resetToken := secure.GenerateToken(32)
	expire := time.Now().Add(24 * time.Hour)
	auth := &models.Auth{Uid: user.Id, Token: resetToken, Type: models.AuthTypeUserFindPwd, ExpiryDate: expire}

	_, err = storage.AddNewAuth(auth)
	if err != nil {
		return errors.New("Add Auth Err:" + err.Error())
	}

	mailer, err := mail.NewServiceMailer()
	if err != nil {
		log.Error("Send Find Password Email Err:", err)
		return err
	}

	resetUrl := beego.AppConfig.String("serverbaseurl") + "/uauth/user/fpwd/email?reset=" + resetToken + "&redirect=" + redirect
	body := fmt.Sprintf(`	尊敬的 %s 您好！
<br>
点击 <a href="%s">链接</a> 重置您账户的密码！
<br>
为保障您的帐号安全，请在24小时内点击该链接，您也可以将链接复制到浏览器地址栏访问。若不是您本人所为，请忽略本邮件，由此给您带来的不便请谅解。
<br>
<br>
本邮件由系统自动发出，请勿直接回复！
<br>
<br>
`, user.Email, resetUrl)
	err = mailer.SendMail(user.Email, "uAuth", "重置密码", "html", body)
	if err != nil {
		log.Error("Send Find Password Email Err:", err)
	}

	return err
}
