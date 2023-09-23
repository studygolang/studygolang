// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

// 动态（go动态；本站动态等）
type Dynamic struct {
	Id      int       `json:"id" xorm:"pk autoincr"`
	Content string    `json:"content"`
	Dmtype  int       `json:"dmtype"`
	Url     string    `json:"url"`
	Seq     int       `json:"seq"`
	Ctime   time.Time `json:"ctime" xorm:"created"`
}
