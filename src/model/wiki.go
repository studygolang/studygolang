// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

type Wiki struct {
	Id      int       `json:"id" xorm:"pk autoincr"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Uri     string    `json:"uri"`
	Uid     int       `json:"uid"`
	Cuid    string    `json:"cuid"`
	Viewnum int       `json:"viewnum"`
	Tags    string    `json:"tags"`
	Ctime   OftenTime `json:"ctime" xorm:"created"`
	Mtime   time.Time `json:"mtime" xorm:"<-"`

	Users map[int]*User `xorm:"-"`
}

func (this *Wiki) BeforeInsert() {
	if this.Tags == "" {
		this.Tags = AutoTag(this.Title, this.Content, 4)
	}
}
