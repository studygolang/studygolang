// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"time"
)

type GCTTUser struct {
	Id        int `xorm:"pk autoincr"`
	Username  string
	Avatar    string
	Uid       int
	JoinedAt  int64
	LastAt    int64
	Num       int
	AvgTime   int
	CreatedAt time.Time `xorm:"<-"`
}

func (*GCTTUser) TableName() string {
	return "gctt_user"
}

type GCTTGit struct {
	Id            int `xorm:"pk autoincr"`
	Username      string
	Title         string
	TranslatingAt int64
	TranslatedAt  int64
	CreatedAt     time.Time `xorm:"<-"`
}

func (*GCTTGit) TableName() string {
	return "gctt_git"
}
