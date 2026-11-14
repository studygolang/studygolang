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

	Users map[int]*User `xorm:"-" json:"-"`

	// 公开投影字段（API 边界填充，避免直接序列化 Users 泄露 Email/Balance 等）
	Author       string            `json:"author" xorm:"-"`
	Contributors []WikiContributor `json:"contributors" xorm:"-"`
}

// WikiContributor 公开贡献者信息，仅含展示所需字段
type WikiContributor struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (this *Wiki) BeforeInsert() {
	if this.Tags == "" {
		this.Tags = AutoTag(this.Title, this.Content, 4)
	}
}
