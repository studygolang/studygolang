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

var BalanceTypeMap = map[int]string{
	MissionTypeLogin:    "每日登录奖励",
	MissionTypeInitial:  "初始资本",
	MissionTypeShare:    "分享获得",
	MissionTypeAdd:      "充值获得",
	MissionTypeReply:    "创建回复",
	MissionTypeTopic:    "创建主题",
	MissionTypeArticle:  "发表文章",
	MissionTypeResource: "分享资源",
	MissionTypeWiki:     "创建WIKI",
	MissionTypeProject:  "发布项目",
	MissionTypeBook:     "分享图书",
	MissionTypeAppend:   "增加附言",
	MissionTypeTop:      "置顶",
	MissionTypeModify:   "修改",
	MissionTypeReplied:  "回复收益",
	MissionTypeAward:    "额外赠予",
	MissionTypeActive:   "活跃奖励",
	MissionTypeGift:     "兑换物品",
	MissionTypePunish:   "处罚",
	MissionTypeSpam:     "Spam",
}

type UserBalanceDetail struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Uid       int       `json:"uid"`
	Type      int       `json:"type"`
	Num       int       `json:"num"`
	Balance   int       `json:"balance"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`

	TypeShow string `json:"type_show" xorm:"-"`
}

func (this *UserBalanceDetail) AfterSet(name string, cell xorm.Cell) {
	if name == "type" {
		this.TypeShow = BalanceTypeMap[this.Type]
	}
}

type UserRecharge struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Uid       int       `json:"uid"`
	Amount    int       `json:"amount"`
	Channel   string    `json:"channel"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}
