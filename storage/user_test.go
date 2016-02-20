package storage

import (
	"testing"
	"uauth/models"
	"uauth/mtest"
)

const (
	existUserEmail  = "work_test_a@163.com"
	defaultPassword = "123456"
)

func TestCreateUser(t *testing.T) {
	mtest.RegisterMySQL()
	user := &models.User{Email: existUserEmail, Password: defaultPassword}
	created, err := CreateUser(user)
	if !created {
		t.Error("ERR:", err)
	}

	user = &models.User{Email: existUserEmail, Password: defaultPassword}
	created, err = CreateUser(user)
	if created {
		t.Error("User %s Allready Exist, must not be created twice!", existUserEmail)
	}

	// delete user
	deleted, err := DeleteUser(user)
	if !deleted {
		t.Error("User should be deleted, ERR:", err)
	}
}
