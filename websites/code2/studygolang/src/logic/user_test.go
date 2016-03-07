package logic_test

import (
	. "logic"
	"testing"
)

func TestFindUserInfos(t *testing.T) {
	usersMap := DefaultUser.FindUserInfos(nil, []int{1, 2, 3})
	t.Fatal(usersMap)
}
