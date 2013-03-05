package model_test

import (
	"encoding/json"
	. "model"
	"testing"
)

func TestNewUserLogin(t *testing.T) {
	userLogin := NewUserLogin()
	userData := `{"uid":"1111","username":"poalris","email":"studygolang@gmail.com","passwd":"123456"}`
	json.Unmarshal([]byte(userData), userLogin)
	// err := userLogin.Find()
	affectedNum, err := userLogin.Insert()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(affectedNum)
}
