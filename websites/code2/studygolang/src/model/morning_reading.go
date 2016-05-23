// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

const (
	RtypeGo   = iota // Go技术晨读
	RtypeComp        // 综合技术晨读
)

// 技术晨读
type MorningReading struct {
	Id       int       `json:"id" xorm:"pk autoincr"`
	Content  string    `json:"content"`
	Rtype    int       `json:"rtype"`
	Inner    int       `json:"inner"`
	Url      string    `json:"url"`
	Moreurls string    `json:"moreurls"`
	Username string    `json:"username"`
	Clicknum int       `json:"clicknum,omitempty"`
	Ctime    OftenTime `json:"ctime" xorm:"<-"`

	// 晨读日期，从 ctime 中提取
	Rdate string `json:"rdate,omitempty" xorm:"-"`

	Urls []string `json:"urls" xorm:"-"`
}

func (this *MorningReading) AfterSet(name string, cell xorm.Cell) {
	switch name {
	case "ctime":
		this.Rdate = time.Time(this.Ctime).Format("2006-01-02")
	case "moreurls":
		if this.Moreurls != "" {
			this.Urls = strings.Split(this.Moreurls, ",")
		}
	}
}
