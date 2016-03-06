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
	// 读取 Post 内容
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

	durationInt, err := strconv.Atoi(duration)
	if err != nil {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("duration must type of int")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if email == "" || pwd == "" {
		resp.SetStatus(http.StatusBadRequest)
		resp.SetMessage("miss Email or Password")
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user := &models.User{Email: email, Password: pwd}

	// 检查邮箱、密码
	err = storage.CheckEmailPwd(user)
	if err == nil {
		if user.Group == models.UserGroupNoActived {
			resp.SetStatus(StatusUserNotActivated)
			resp.SetMessage(StatusText(StatusUserNotActivated))
			this.Data["json"] = resp
			this.ServeJSON()
			return
		}

		token := secure.GenerateToken(32)
		expiry := time.Now().Add(time.Duration(durationInt) * time.Hour)
		auth := &models.Auth{Uid: user.Id, Token: token, Server: web, Status: "ok", Type: models.AuthTypeUserSignIn, ExpiryDate: expiry}
		_, err = storage.AddNewAuth(auth)
		if err != nil {
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
		resp.SetStatus(StatusWrongUserNameOrPwd)
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
