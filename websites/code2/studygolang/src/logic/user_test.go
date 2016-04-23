package logic_test

import (
	. "logic"
	"testing"
)

func TestFindUserInfos(t *testing.T) {
	usersMap := DefaultUser.FindUserInfos(nil, []int64{1, 2, 3})
	if len(usersMap) == 0 {
		t.Fatal(usersMap)
	}
}
