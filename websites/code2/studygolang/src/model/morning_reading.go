// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package model

// 资源信息
type MorningReading struct {
	Id       int    `json:"id" xorm:"pk autoincr"`
	Content  string `json:"content"`
	Rtype    int    `json:"rtype"`
	Inner    int    `json:"inner"`
	Url      string `json:"url"`
	Moreurls string `json:"moreurls"`
	Username string `json:"username"`
	Clicknum int    `json:"clicknum,omitempty"`
	Ctime    string `json:"ctime,omitempty" xorm:"<-"`
}
