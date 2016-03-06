package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"uauth/models"
	"uauth/storage"
)

type UserActiveHandler struct {
	BaseController
}

func (this *UserActiveHandler) ActiveFromEmail() {
	activeToken := this.GetString("active", "")
	auth, err := storage.FindAuthByToken(activeToken)

	// 没有此 token
	if err != nil {
		this.Ctx.WriteString("404! token not found!")
		return
	}

	redirect := auth.Redirect
	// token 过期
	if time.Now().After(auth.ExpiryDate) {
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusForbidden, "token is expired")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("403! token is expired!")
		}
		return
	}

	uid := auth.Uid

	user := &models.User{Id: uid, Group: models.UserGroupCommon}

	err = storage.UpdateUser(user, "Group")
	if err != nil {
		log.Error("Update User Err:", err)
		if redirect != "" {
			redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusInternalServerError, "active failed")
			this.Redirect(redirect, 302)
		} else {
			this.Ctx.WriteString("500! active failed!")
		}
		return
	}

	if redirect != "" {
		redirect = fmt.Sprintf(redirect+"?status=%d&msg=%s", http.StatusOK, "user actived success")
		this.Redirect(redirect, 302)
	} else {
		this.Ctx.WriteString("200! user actived success")
	}

}

// ResendActivateEmail 重新发送激活邮件
func (this *UserActiveHandler) ResendActivateEmail() {
	resp := make(Response)

	body := this.Ctx.Request.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("Read Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	rqsBd := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &rqsBd)
	if err != nil {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	email := rqsBd["Email"]
	redirect := rqsBd["Redirect"]
	if email == "" {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("miss Email")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user, err := storage.FindUserByEmail(email)
	if err != nil {
		resp.SetStatus(http.StatusNotFound)
		resp.SetMessage("User Not Exist")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	err = sendActiveEmail(user, redirect)
	if err != nil {
		resp.SetStatus(StatusSendEmailFailed)
		resp.SetMessage("send activate user email failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp.SetStatus(http.StatusOK)
	resp.SetMessage("send activate user email success")
	this.Data["json"] = resp
	this.ServeJSON()
}
