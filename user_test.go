package main

import (
	"testing"
)

func TestNewUserNotEmptyData(t *testing.T) {

	cfg := NewConfig()
	u := NewUser(cfg)

	if u.login == "" || u.password == "" || string(u.passwordHash) == "" {
		t.Errorf("User struct should not have an empty values: got %+v", u)
	}
}
