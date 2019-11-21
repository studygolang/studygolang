// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

type ViewSource struct {
	Id        int `xorm:"pk autoincr"`
	Objid     int
	Objtype   int
	Google    int
	Baidu     int
	Bing      int
	Sogou     int
	So        int
	Other     int
	UpdatedAt OftenTime `xorm:"<-"`
}
