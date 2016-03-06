package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"uauth/models"
	"uauth/storage"
)

type UserDataHandler struct {
	BaseController
}

func (this *UserDataHandler) FindUserWithAuthToken() {
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

	rspBd := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &rspBd)
	if err != nil {
		log.Error("Parse Request Body Err:", err)
		resp.SetStatus(StatusUnprocessableEntity)
		resp.SetMessage("Parse Body Failed")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	authToken := rspBd["AuthToken"]
	auth, err := storage.FindAuthByToken(authToken)
	if err != nil {
		resp.SetStatus(http.StatusForbidden)
		resp.SetMessage("this token is not available")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if time.Now().After(auth.ExpiryDate) {
		resp.SetStatus(http.StatusForbidden)
		resp.SetMessage("this token is expired")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if auth.Type != models.AuthTypeUserSignIn {
		resp.SetStatus(http.StatusForbidden)
		resp.SetMessage("Token miss type")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	uid := auth.Uid
	user, err := storage.FindUserById(uid)
	if err != nil {
		resp.SetStatus(http.StatusNotFound)
		resp.SetMessage("can not find user by this token")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if user.Group == models.UserGroupNoActived {
		resp.SetStatus(StatusUserNotActivated)
		resp.SetMessage("user not activated")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp.SetStatus(http.StatusOK)
	resp.SetMessage("ok")
	resp.SetData("User", user)
	this.Data["json"] = resp
	this.ServeJSON()
}
