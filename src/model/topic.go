// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	FlagNoAudit = iota
	FlagNormal
	FlagAuditDelete
	FlagUserDelete
)

const (
	// 最多附言条数
	AppendMaxNum = 3
)

const (
	PermissionPublic = iota // 公开
	PermissionLogin         // 登录可见
	PermissionFollow        // 关注可见（暂未实现）
)

// 社区主题信息
type Topic struct {
	Tid           int       `xorm:"pk autoincr" json:"tid"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Nid           int       `json:"nid"`
	Uid           int       `json:"uid"`
	Flag          uint8     `json:"flag"`
	Lastreplyuid  int       `json:"lastreplyuid"`
	Lastreplytime OftenTime `json:"lastreplytime"`
	EditorUid     int       `json:"editor_uid"`
	Top           uint8     `json:"top"`
	TopTime       int64     `json:"top_time"`
	Tags          string    `json:"tags"`
	Permission    int       `json:"permission"`
	Ctime         OftenTime `json:"ctime" xorm:"created"`
	Mtime         OftenTime `json:"mtime" xorm:"<-"`

	// 为了方便，加上Node（节点名称，数据表没有）
	Node string `xorm:"-"`
	// 排行榜阅读量
	RankView int `json:"rank_view" xorm:"-"`
}

func (*Topic) TableName() string {
	return "topics"
}

func (this *Topic) BeforeInsert() {
	if this.Tags == "" {
		this.Tags = AutoTag(this.Title, this.Content, 4)
	}
}

// 社区主题扩展（计数）信息
type TopicEx struct {
	Tid   int       `json:"-"`
	View  int       `json:"view"`
	Reply int       `json:"reply"`
	Like  int       `json:"like"`
	Mtime time.Time `json:"mtime" xorm:"<-"`
}

func (*TopicEx) TableName() string {
	return "topics_ex"
}

// 社区主题扩展（计数）信息，用于 incr 更新
type TopicUpEx struct {
	Tid   int       `json:"-" xorm:"pk"`
	View  int       `json:"view"`
	Reply int       `json:"reply"`
	Like  int       `json:"like"`
	Mtime time.Time `json:"mtime" xorm:"<-"`
}

func (*TopicUpEx) TableName() string {
	return "topics_ex"
}

type TopicInfo struct {
	Topic   `xorm:"extends"`
	TopicEx `xorm:"extends"`
}

func (*TopicInfo) TableName() string {
	return "topics"
}

type TopicAppend struct {
	Id        int `xorm:"pk autoincr"`
	Tid       int
	Content   string
	CreatedAt OftenTime `xorm:"<-"`
}

// 社区主题节点信息
type TopicNode struct {
	Nid       int       `json:"nid" xorm:"pk autoincr"`
	Parent    int       `json:"parent"`
	Logo      string    `json:"logo"`
	Name      string    `json:"name"`
	Ename     string    `json:"ename"`
	Seq       int       `json:"seq"`
	Intro     string    `json:"intro"`
	ShowIndex bool      `json:"show_index"`
	Ctime     time.Time `json:"ctime" xorm:"<-"`

	Level int `json:"-" xorm:"-"`
}

func (*TopicNode) TableName() string {
	return "topics_node"
}

// 推荐节点
type RecommendNode struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	Parent    int       `json:"parent"`
	Nid       int       `json:"nid"`
	Seq       int       `json:"seq"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}

type NodeInfo struct {
	RecommendNode `xorm:"extends"`
	TopicNode     `xorm:"extends"`
}

func (*NodeInfo) TableName() string {
	return "recommend_node"
}
