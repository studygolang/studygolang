package model_test

import (
	. "model"
	"testing"
	//"encoding/json"
)

func TestNewUserLogin(t *testing.T) {
	userLogin := NewUserLogin()
	// userData := `{"uid":123234,"username":"poalris","email":"studygolang@gmail.com","passwd":"123456"}`
	// json.Unmarshal([]byte(userData), userLogin)
	err := userLogin.Find()
	// affectedNum, err := userLogin.Insert()
	if err != nil {
	    t.Fatal(err)
	}
	t.Error(userLogin.Uid)
}