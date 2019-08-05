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
	Id          int `json:"id" xorm:"pk autoincr"`
	Name        string
	Description string
	Price       int
	TotalNum    int
	RemainNum   int
	ExpireTime  time.Time `xorm:"int"`
	Supplier    string
	BuyLimit    int
	Typ         int
	State       int
	CreatedAt   OftenTime `xorm:"<-"`

	TypShow string `xorm:"-"`
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
	Id         int `json:"id" xorm:"pk autoincr"`
	GiftId     int
	Uid        int
	Remark     string
	ExpireTime time.Time `xorm:"int"`
	CreatedAt  OftenTime `xorm:"<-"`
}
