// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package model

import "time"

const (
	FlagNoAudit = iota
	FlagNormal
	FlagAuditDelete
	FlagUserDelete
)

// 社区主题信息
type Topic struct {
	Tid           int       `gorm:"primary_key" json:"tid"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Nid           int       `json:"nid"`
	Uid           int       `json:"uid"`
	Flag          uint8     `json:"flag"`
	Lastreplyuid  int       `json:"lastreplyuid"`
	Lastreplytime time.Time `json:"lastreplytime"`
	EditorUid     int       `json:"editor_uid"`
	Top           bool      `json:"istop"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

func (*Topic) TableName() string {
	return "topics"
}

// 社区主题扩展（计数）信息
type TopicEx struct {
	Tid   int       `gorm:"primary_key" json:"tid"`
	View  int       `json:"view"`
	Reply int       `json:"reply"`
	Like  int       `json:"like"`
	Mtime time.Time `json:"mtime"`
}

func (*TopicEx) TableName() string {
	return "topics_ex"
}

// 社区主题节点信息
type TopicNode struct {
	Nid    int       `json:"nid" gorm:"primary_key"`
	Parent int       `json:"parent"`
	Name   string    `json:"name"`
	Intro  string    `json:"intro"`
	Ctime  time.Time `json:"ctime"`
}

func (*TopicNode) TableName() string {
	return "topics_node"
}
