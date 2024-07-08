package controller

import "testing"

func TestUser(t *testing.T) {
	controller := resetDb(t)
	_, e := controller.CreateUser("", "123", false)
	if e == nil {
		t.Error("blank username did not fail")
	}
	_, e = controller.CreateUser("cool_username_123", "super_secure_password", false)
	if e != nil {
		t.Error("user creation failed")
	}

	goodId, loginErr := controller.LoginUser("cool_username_123", "super_secure_password")
	if goodId < 0 || loginErr != nil {
		t.Error("failed to login")
	}
	id, loginErr := controller.LoginUser("cool_username_123", "supre_secure_password")
	if id > 0 || loginErr == nil {
		t.Error("successful login with incorrect password")
	}
	if loginErr == nil {
		t.Error("no error on bad login")
	}
	id, loginErr = controller.LoginUser("dne", "super_secure_password")
	if id > 0 || loginErr == nil {
		t.Error("successful login with non-existent user")
	}
	delErr := controller.DeleteUser(goodId)
	if delErr != nil {
		t.Errorf("failed to delete user [%d] %s\n", goodId, delErr)
	}
	delErr = controller.DeleteUser(goodId)
	if delErr == nil {
		t.Error("no error on double delete")
	}
	delErr = controller.DeleteUser(9999)
	if delErr == nil {
		t.Error("no error on delete non existent")
	}
}
