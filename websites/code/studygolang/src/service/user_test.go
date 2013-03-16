package service_test

import (
	. "service"
	"testing"
)

func TestFindUsers(t *testing.T) {
	userList, err := FindUsers()
	if err != nil && len(userList) == 0 {
		t.Fatal(err)
	}
	t.Log(len(userList))
	for k, tmpUser := range userList {
		t.Log(k, tmpUser)
		t.Log("===")
	}
}
