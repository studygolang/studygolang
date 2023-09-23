// Copyright 2022 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"strconv"
	"time"

	"xorm.io/xorm"
)

// Go 面试题
type InterviewQuestion struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Sn        int64     `json:"sn"`
	ShowSn    string    `json:"show_sn" xorm:"-"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Level     int       `json:"level"`
	Viewnum   int       `json:"viewnum"`
	Cmtnum    int       `json:"cmtnum"`
	Likenum   int       `json:"likenum"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}

func (iq *InterviewQuestion) AfterSet(name string, cell xorm.Cell) {
	if name == "sn" {
		iq.ShowSn = strconv.FormatInt(iq.Sn, 32)
	}
}

func (iq *InterviewQuestion) AfterInsert() {
	iq.ShowSn = strconv.FormatInt(iq.Sn, 32)
}
