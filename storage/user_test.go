package storage

import (
	"testing"
	"uauth/models"
	"uauth/mtest"
)

const (
	signUpUserEmail = "work_test_a@163.com"
	signInUserEmail = "work_test_b@163.com"
	defaultPassword = "123456"
)

func TestCreateUser(t *testing.T) {
	mtest.RegisterMySQL()
	user := &models.User{Email: signUpUserEmail, Password: defaultPassword}
	created, err := CreateUser(user)
	if !created {
		t.Error("ERR:", err)
	}

	user = &models.User{Email: signUpUserEmail, Password: defaultPassword}
	created, err = CreateUser(user)
	if created {
		t.Error("User %s Allready Exist, must not be created twice", signUpUserEmail)
	}

	// delete user
	deleted, err := DeleteUser(user)
	if !deleted {
		t.Error("User should be deleted, ERR:", err)
	}
}

func TestCheckEmailPwd(t *testing.T) {
	user := &models.User{Email: signInUserEmail, Password: defaultPassword}
	err := CheckEmailPwd(user)
	if err != nil {
		t.Error("TestCheckEmailPwd Fail, err expected nil, actual:", err)
	} else if user.Id == 0 {
		t.Error("TestCheckEmailPwd Fail, user.Id must not 0")
	}

	err = CheckEmailPwd(&models.User{Email: signUpUserEmail, Password: defaultPassword})
	if err == nil {
		t.Error("TestCheckEmailPwd Fail, err should be nil")
	}
}
