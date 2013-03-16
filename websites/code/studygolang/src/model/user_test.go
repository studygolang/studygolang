// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model_test

import (
	"encoding/json"
	. "model"
	"testing"
)

func TestNewUserLogin(t *testing.T) {
	user := NewUser()
	userList, err := user.FindAll()
	for _, tmpUser := range userList {
		t.Log(tmpUser.Name)
		t.Log("===")
	}
	if err == nil {
		t.Fatal(err)
	}
}

func testInsert(t *testing.T) {
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
