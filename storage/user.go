package storage

import (
	"encoding/hex"
	"uauth/models"

	"crypto/md5"
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

// CreateUser create a new user in mysql, if user is allready exist in table `user`,
// or some errors occurred during creating user, return error
func CreateUser(user *models.User) (created bool, err error) {
	// if user is nil, return error
	if user == nil {
		return false, errors.New("0-user must not nil!")
	}
	// if email or pwd is empty, return error
	if user.Email == "" || user.Password == "" {
		return false, errors.New("1-user Email or Password must not empty!")
	}

	encryptoPass, err := doMd5(user.Password)
	if err != nil {
		log.Error("MD5 ERR:", err)
		return false, errors.New("2-MD5 ERR:" + err.Error())
	}

	// init default UserName DisplayName by Email
	user.Password = encryptoPass
	user.UserName, user.DisplayName, user.Created, user.Activated, user.Logged = user.Email, user.Email, time.Now(), time.Now(), time.Now()

	// save user to mysql table `user`
	o := orm.NewOrm()
	created, id, err := o.ReadOrCreate(user, "Email")

	if err != nil {
		return false, err
	}

	if created {
		user.Id = id
		log.Info("Create User Success:", user.String())
		return true, nil
	}
	err = errors.New("3-User Allready Exist!")
	return false, err
}

// DeleteUser delete user in mysql table `user`, user.Id must be setted
func DeleteUser(user *models.User) (deleted bool, err error) {
	if user == nil {
		return false, errors.New("user must not nil!")
	}

	if user.Id == 0 {
		return false, errors.New("user Id must not 0")
	}

	o := orm.NewOrm()
	num, err := o.Delete(user)
	if err != nil {
		return false, err
	}

	if num != 1 {
		log.Warn("Delete Multi User: ", num)
	}
	return true, err
}

// doMd5 return md5-ed data
func doMd5(data string) (encrypto string, err error) {
	md5Hash := md5.New()
	_, err = md5Hash.Write([]byte(data))
	if err != nil {
		return "", err
	}
	hashedBytes := md5Hash.Sum(nil)
	hashedData := hex.EncodeToString(hashedBytes)
	return hashedData, nil
}
