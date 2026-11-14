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
	GiftStateOnline  = 1
	GiftStateExpired = 3

	GiftTypRedeem   = 0
	GiftTypDiscount = 1
)

var GiftTypeMap = map[int]string{
	GiftTypRedeem:   "兑换码",
	GiftTypDiscount: "折扣",
}

type Gift struct {
	Id          int       `json:"id" xorm:"pk autoincr"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	TotalNum    int       `json:"total_num"`
	RemainNum   int       `json:"remain_num"`
	ExpireTime  time.Time `json:"expire_time" xorm:"int"`
	Supplier    string    `json:"supplier"`
	BuyLimit    int       `json:"buy_limit"`
	Typ         int       `json:"typ"`
	State       int       `json:"state"`
	CreatedAt   OftenTime `json:"created_at" xorm:"<-"`

	TypShow string `json:"typ_show" xorm:"-"`

	// CanExchange 由 GiftLogic.UserCanExchange 计算，未登录时恒为 false
	CanExchange bool `json:"can_exchange" xorm:"-"`
}

func (this *Gift) AfterSet(name string, cell xorm.Cell) {
	if name == "typ" {
		this.TypShow = GiftTypeMap[this.Typ]
	}
}

type GiftRedeem struct {
	Id        int `json:"id" xorm:"pk autoincr"`
	GiftId    int
	Code      string
	Exchange  int
	Uid       int
	UpdatedAt OftenTime `xorm:"<-"`
}

type UserExchangeRecord struct {
	Id         int       `json:"id" xorm:"pk autoincr"`
	GiftId     int       `json:"gift_id"`
	Uid        int       `json:"uid"`
	Remark     string    `json:"remark"`
	ExpireTime time.Time `json:"expire_time" xorm:"int"`
	CreatedAt  OftenTime `json:"created_at" xorm:"<-"`
}
