// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

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
