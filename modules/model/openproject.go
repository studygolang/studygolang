// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"time"

	"github.com/go-xorm/xorm"
)

const (
	ProjectStatusNew     = 0
	ProjectStatusOnline  = 1
	ProjectStatusOffline = 2
)

// 开源项目信息
type OpenProject struct {
	Id            int       `json:"id" xorm:"pk autoincr"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Uri           string    `json:"uri"`
	Home          string    `json:"home"`
	Doc           string    `json:"doc"`
	Download      string    `json:"download"`
	Src           string    `json:"src"`
	Logo          string    `json:"logo"`
	Desc          string    `json:"desc"`
	Repo          string    `json:"repo"`
	Author        string    `json:"author"`
	Licence       string    `json:"licence"`
	Lang          string    `json:"lang"`
	Os            string    `json:"os"`
	Tags          string    `json:"tags"`
	Username      string    `json:"username,omitempty"`
	Viewnum       int       `json:"viewnum,omitempty"`
	Cmtnum        int       `json:"cmtnum,omitempty"`
	Likenum       int       `json:"likenum,omitempty"`
	Lastreplyuid  int       `json:"lastreplyuid"`
	Lastreplytime OftenTime `json:"lastreplytime"`
	Status        int       `json:"status"`
	Ctime         OftenTime `json:"ctime,omitempty" xorm:"created"`
	Mtime         OftenTime `json:"mtime,omitempty" xorm:"<-"`

	User *User `json:"user" xorm:"-"`
	// 排行榜阅读量
	RankView      int   `json:"rank_view" xorm:"-"`
	LastReplyUser *User `json:"last_reply_user" xorm:"-"`
}

func (this *OpenProject) BeforeInsert() {
	if this.Tags == "" {
		this.Tags = AutoTag(this.Name+this.Category, this.Desc, 4)
	}

	this.Lastreplytime = NewOftenTime()
}

func (this *OpenProject) AfterInsert() {
	go func() {
		// AfterInsert 时，自增 ID 还未赋值，这里 sleep 一会，确保自增 ID 有值
		for {
			if this.Id > 0 {
				PublishFeed(this, nil)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (this *OpenProject) AfterSet(name string, cell xorm.Cell) {
	if name == "logo" && this.Logo == "" {
		this.Logo = WebsiteSetting.ProjectDfLogo
	}
}
