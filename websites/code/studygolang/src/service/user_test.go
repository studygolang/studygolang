// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

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
