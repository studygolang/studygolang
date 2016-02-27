package logic_test

import (
	"logic"
	"testing"

	"github.com/polaris1119/dbutil"
)

func init() {
	dbutil.InitDB("root:@tcp(localhost:3306)/studygolang?charset=utf8")

}

func TestFindUserInfos(t *testing.T) {
	usersMap := logic.DefaultUserLogic.FindUserInfos(nil, []int{1, 2, 3})
	t.Fatal(usersMap)
}
