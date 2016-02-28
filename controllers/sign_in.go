package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"uauth/models"
	"uauth/storage"
)

type SignInHandler struct {
	BaseController
}

func (this *SignInHandler) SignIn() {
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

	err = storage.CheckEmailPwd(user)
	if err == nil {
		log.Info("Sign In Success:", user.String())
		if user.Group == models.UserGroupNoActived {
			resp.SetStatus(http.StatusUnauthorized)
			resp.SetMessage("user not actived")
			this.Data["json"] = resp
			this.ServeJSON()
			return
		}
		resp.SetStatus(http.StatusOK)
		resp.SetMessage("OK")
		resp.SetData(user)
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if err == storage.ErrNoRows {
		log.Info("Sign In Err: Email or Password is wrong")
		resp.SetStatus(StatusUserExist)
		resp.SetMessage("Email or Password is wrong")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	log.Error("Sign In Err: ", err)
	resp.SetStatus(StatusUkownError)
	resp.SetMessage("Sign In Failed")
	this.Data["json"] = resp
	this.ServeJSON()
}
