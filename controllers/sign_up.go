package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"uauth/models"
	"uauth/storage"
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
