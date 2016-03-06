package storage

import (
	"uauth/models"

	"errors"
	"github.com/astaxie/beego/orm"
	"time"
	"uauth/tools/secure"
)

// CreateUser create a new user in mysql, if user is allready exist in table `user`,
// or some errors occurred during creating user, return error
func CreateUser(user *models.User) (created bool, err error) {
	// if user is nil, return error
	if user == nil {
		return false, errors.New("0-user must not nil")
	}
	// if email or pwd is empty, return error
	if user.Email == "" || user.Password == "" {
		return false, errors.New("1-user Email or Password must not empty")
	}

	encryptoPass, err := secure.DoMd5(user.Password)
	if err != nil {
		log.Error("MD5 ERR:", err)
		return false, errors.New("2-MD5 ERR:" + err.Error())
	}

	// init default UserName DisplayName by Email
	user.Password = encryptoPass
	user.UserName, user.DisplayName, user.Created, user.Activated, user.Logged, user.Group =
		user.Email, user.Email, time.Now(), time.Now(), time.Now(), models.UserGroupNoActived

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

	return false, ErrRowExist
}

// DeleteUser delete user in mysql table `user`, user.Id must be setted
func DeleteUser(user *models.User) (deleted bool, err error) {
	if user == nil {
		return false, errors.New("user must not nil")
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

// FindUserById 通过 User:Id 查询 User
func FindUserById(uid int64) (user *models.User, err error) {
	o := orm.NewOrm()
	user = &models.User{Id: uid}
	err = o.Read(user, "Id")
	if err != nil {
		return nil, ErrNoRows
	}
	return user, nil
}

// FindUserByEmail 通过 User:Email 查询 User
func FindUserByEmail(email string) (user *models.User, err error) {
	o := orm.NewOrm()
	user = &models.User{Email: email}
	err = o.Read(user, "Email")
	if err != nil {
		return nil, ErrNoRows
	}

	return user, nil
}

func UpdateUser(user *models.User, columns ...string) (err error) {
	o := orm.NewOrm()

	if user != nil && user.Id == 0 {
		return errors.New("user is nil or miss user pk:Id")
	}

	_, err = o.Update(user, columns...)

	return err
}

// CheckEmailPwd  check user email and password in mysql table user
func CheckEmailPwd(user *models.User) (err error) {
	if user == nil {
		return errors.New("0-user must not nil")
	}

	if user.Email == "" || user.Password == "" {
		return errors.New("1-user.Eamil or user.Password must not empty")
	}

	encryptoPass, err := secure.DoMd5(user.Password)
	if err != nil {
		log.Error("MD5 ERR:", err)
		return errors.New("2-MD5 ERR:" + err.Error())
	}
	user.Password = encryptoPass
	o := orm.NewOrm()
	err = o.Read(user, "Email", "Password")

	if err == orm.ErrNoRows {
		return ErrNoRows
	}
	if err != nil {
		return errors.New("4-Read Table User Err:" + err.Error())
	}
	return nil
}
