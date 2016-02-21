package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
		resp.SetStatus(http.StatusOK)
		resp.SetMessage("OK")
		resp.SetData(user)
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	errMsg := strings.Split(err.Error(), "-")
	errCode := errMsg[0]
	// errCode == 3, user exist
	if errCode == "3" {
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
