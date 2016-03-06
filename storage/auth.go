package storage

import (
	"uauth/models"

	"github.com/astaxie/beego/orm"
)

// AddNewAuth 保存一个新的 Auth
func AddNewAuth(auth *models.Auth) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(auth)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// FindAuthByToken 通过 Token 查询 Auth
func FindAuthByToken(token string) (auth *models.Auth, err error) {
	auth = &models.Auth{Token: token}
	o := orm.NewOrm()
	err = o.QueryTable(auth).Filter("Token", token).One(auth)

	return auth, err
}
