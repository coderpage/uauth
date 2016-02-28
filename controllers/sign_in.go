package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"uauth/models"
	"uauth/storage"
	"uauth/tools/secure"
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

	email := rspBd["Email"]
	pwd := rspBd["Password"]
	web := rspBd["Web"]
	duration := rspBd["Duration"]

	user := &models.User{Email: email, Password: pwd}
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

		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			resp.SetStatus(StatusUnprocessableEntity)
			resp.SetMessage("duration must type of int")
			this.Data["json"] = resp
			this.ServeJSON()
			return
		}

		token := secure.GenerateToken(32)
		expiry := time.Now().Add(time.Duration(durationInt) * time.Hour)
		auth := &models.Auth{Uid: user.Id, Token: token, Server: web, Status: "ok", Type: models.AuthTypeUserSignIn, ExpiryDate: expiry}
		_, err = storage.AddNewAuth(auth)
		if err != nil {
			log.Error("save auth failed:", err)
			resp.SetStatus(http.StatusInternalServerError)
			resp.SetMessage("save token failed")
			this.Data["json"] = resp
			this.ServeJSON()
			return
		}

		resp.SetStatus(http.StatusOK)
		resp.SetMessage("OK")
		resp.SetData("User", user)
		resp.SetData("AuthToken", token)
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
