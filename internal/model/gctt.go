// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"time"

	"xorm.io/xorm"
)

const (
	GCTTRoleTranslator = iota
	GCTTRoleLeader
	GCTTRoleSelecter // 选题
	GCTTRoleChecker  // 校对
	GCTTRoleCore     // 核心成员
)

const (
	IssueOpened = iota
	IssueClosed
)

const (
	LabelUnClaim = "待认领"
	LabelClaimed = "已认领"
)

var roleMap = map[int]string{
	GCTTRoleTranslator: "译者",
	GCTTRoleLeader:     "组长",
	GCTTRoleSelecter:   "选题",
	GCTTRoleChecker:    "校对",
	GCTTRoleCore:       "核心成员",
}

var faMap = map[int]string{
	GCTTRoleTranslator: "fa-user",
	GCTTRoleLeader:     "fa-graduation-cap",
	GCTTRoleSelecter:   "fa-user-circle",
	GCTTRoleChecker:    "fa-user-secret",
	GCTTRoleCore:       "fa-heart",
}

type GCTTUser struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Uid       int       `json:"uid"`
	JoinedAt  int64     `json:"joined_at"`
	LastAt    int64     `json:"last_at"`
	Num       int       `json:"num"`
	Words     int       `json:"words"`
	AvgTime   int       `json:"avg_time"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`

	RoleName string `json:"role_name" xorm:"-"`
	Fa       string `json:"fa" xorm:"-"`
}

func (this *GCTTUser) AfterSet(name string, cell xorm.Cell) {
	if name == "role" {
		this.RoleName = roleMap[this.Role]
		this.Fa = faMap[this.Role]
	}
}

func (*GCTTUser) TableName() string {
	return "gctt_user"
}

type GCTTGit struct {
	Id            int `xorm:"pk autoincr"`
	Username      string
	Md5           string
	Title         string
	PR            int `xorm:"pr"`
	TranslatingAt int64
	TranslatedAt  int64
	Words         int
	ArticleId     int
	CreatedAt     time.Time `xorm:"<-"`
}

func (*GCTTGit) TableName() string {
	return "gctt_git"
}

type GCTTIssue struct {
	Id            int       `json:"id" xorm:"pk autoincr"`
	Translator    string    `json:"translator"`
	Email         string    `json:"-"`
	Title         string    `json:"title"`
	TranslatingAt int64     `json:"translating_at"`
	TranslatedAt  int64     `json:"translated_at"`
	Label         string    `json:"label"`
	State         uint8     `json:"state"`
	CreatedAt     time.Time `json:"created_at" xorm:"<-"`
}

func (*GCTTIssue) TableName() string {
	return "gctt_issue"
}

type GCTTTimeLine struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (*GCTTTimeLine) TableName() string {
	return "gctt_timeline"
}
