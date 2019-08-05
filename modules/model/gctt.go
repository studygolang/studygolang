// Copyright 2017 The StudyGolang Authors. All rights reserved.
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
	Id        int `xorm:"pk autoincr"`
	Username  string
	Avatar    string
	Uid       int
	JoinedAt  int64
	LastAt    int64
	Num       int
	Words     int
	AvgTime   int
	Role      int
	CreatedAt time.Time `xorm:"<-"`

	RoleName string `xorm:"-"`
	Fa       string `xorm:"-"`
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
	Id            int `xorm:"pk autoincr"`
	Translator    string
	Email         string
	Title         string
	TranslatingAt int64
	TranslatedAt  int64
	Label         string
	State         uint8
	CreatedAt     time.Time `xorm:"<-"`
}

func (*GCTTIssue) TableName() string {
	return "gctt_issue"
}

type GCTTTimeLine struct {
	Id        int `xorm:"pk autoincr"`
	Content   string
	CreatedAt time.Time
}

func (*GCTTTimeLine) TableName() string {
	return "gctt_timeline"
}
