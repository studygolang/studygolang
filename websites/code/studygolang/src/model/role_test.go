package model_test

import (
	. "model"
	"testing"
)

func TestNewRole(t *testing.T) {
	roleList, err := NewRole().FindAll()
	for _, tmpUser := range roleList {
		t.Log(tmpUser.Roleid)
		t.Log("===")
	}
	if err != nil {
		t.Fatal(err)
	}
}
