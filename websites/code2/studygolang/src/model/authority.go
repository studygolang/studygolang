// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

// 权限信息
type Authority struct {
	Aid    int       `json:"aid" xorm:"pk autoincr"`
	Name   string    `json:"name"`
	Menu1  int       `json:"menu1"`
	Menu2  int       `json:"menu2"`
	Route  string    `json:"route"`
	OpUser string    `json:"op_user"`
	Ctime  OftenTime `json:"ctime" xorm:"created"`
	Mtime  OftenTime `json:"mtime" xorm:"<-"`
}
