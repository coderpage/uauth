package controllers

import (
	"fmt"
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
