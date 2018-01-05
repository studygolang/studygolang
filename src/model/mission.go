// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	MissionTypeLogin   = 1
	MissionTypeInitial = 2
	MissionTypeShare   = 3
	MissionTypeAdd     = 4

	// 回复
	MissionTypeReply = 51
	// 创建
	MissionTypeTopic    = 52
	MissionTypeArticle  = 53
	MissionTypeResource = 54
	MissionTypeWiki     = 55
	MissionTypeProject  = 56
	MissionTypeBook     = 57

	MissionTypeAppend = 60
	// 置顶
	MissionTypeTop = 61

	MissionTypeModify = 65
	// 被回复
	MissionTypeReplied = 70
	// 额外赠予
	MissionTypeAward = 80
	// 活跃奖励
	MissionTypeActive = 81

	// 物品兑换
	MissionTypeGift = 100

	// 管理员操作后处罚
	MissionTypePunish = 120
	// 水
	MissionTypeSpam = 127
)

const (
	InitialMissionId = 1
)

type Mission struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	Fixed     int       `json:"fixed"`
	Min       int       `json:"min"`
	Max       int       `json:"max"`
	Incr      int       `json:"incr"`
	State     int       `json:"state"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}

type UserLoginMission struct {
	Uid       int       `json:"uid" xorm:"pk"`
	Date      int       `json:"date"`
	Award     int       `json:"award"`
	Days      int       `json:"days"`
	TotalDays int       `json:"total_days"`
	UpdatedAt time.Time `json:"updated_at" xorm:"<-"`
}
