// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	FlagCancel = iota
	FlagLike   // 喜欢
	FlagUnlike // 不喜欢（暂时不支持）
)

// 赞（喜欢）
type Like struct {
	Uid     int       `json:"uid"`
	Objid   int       `json:"objid"`
	Objtype int       `json:"objtype"`
	Flag    int       `json:"flag"`
	Ctime   time.Time `json:"ctime" xorm:"<-"`
}

func (*Like) TableName() string {
	return "likes"
}
