// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"time"
)

// Subject 专栏
type Subject struct {
	Id          int       `xorm:"pk autoincr" json:"id"`
	Name        string    `json:"name"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	Uid         int       `json:"uid"`
	Contribute  bool      `json:"contribute"`
	Audit       bool      `json:"audit"`
	ArticleNum  int       `json:"article_num"`
	CreatedAt   OftenTime `json:"created_at" xorm:"created"`
	UpdatedAt   OftenTime `json:"updated_at" xorm:"<-"`

	User *User `json:"user" xorm:"-"`
}

// SubjectAdmin 专栏管理员
type SubjectAdmin struct {
	Id        int       `xorm:"pk autoincr" json:"id"`
	Sid       int       `json:"sid"`
	Uid       int       `json:"uid"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}

const (
	ContributeStateNew = iota
	ContributeStateOnline
	ContributeStateOffline
)

// SubjectArticle 专栏文章
type SubjectArticle struct {
	Id        int       `xorm:"pk autoincr" json:"id"`
	Sid       int       `json:"sid"`
	ArticleId int       `json:"article_id"`
	State     int       `json:"state"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}

// SubjectArticles xorm join 需要
type SubjectArticles struct {
	Article   `xorm:"extends"`
	Sid       int
	CreatedAt time.Time
}

func (*SubjectArticles) TableName() string {
	return "articles"
}

// SubjectFollower 专栏关注者
type SubjectFollower struct {
	Id        int       `xorm:"pk autoincr" json:"id"`
	Sid       int       `json:"sid"`
	Uid       int       `json:"uid"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`

	User    *User  `xorm:"-"`
	TimeAgo string `xorm:"-"`
}
