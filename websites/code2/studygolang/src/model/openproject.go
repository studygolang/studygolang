// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	ProjectStatusNew     = 0
	ProjectStatusOnline  = 1
	ProjectStatusOffline = 2
)

// 开源项目信息
type OpenProject struct {
	Id       int       `json:"id" xorm:"pk autoincr"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Uri      string    `json:"uri"`
	Home     string    `json:"home"`
	Doc      string    `json:"doc"`
	Download string    `json:"download"`
	Src      string    `json:"src"`
	Logo     string    `json:"logo"`
	Desc     string    `json:"desc"`
	Repo     string    `json:"repo"`
	Author   string    `json:"author"`
	Licence  string    `json:"licence"`
	Lang     string    `json:"lang"`
	Os       string    `json:"os"`
	Tags     string    `json:"tags"`
	Username string    `json:"username,omitempty"`
	Viewnum  int       `json:"viewnum,omitempty"`
	Cmtnum   int       `json:"cmtnum,omitempty"`
	Likenum  int       `json:"likenum,omitempty"`
	Status   int       `json:"status"`
	Ctime    OftenTime `json:"ctime,omitempty" xorm:"created"`
	Mtime    time.Time `json:"mtime,omitempty" xorm:"<-"`
}
