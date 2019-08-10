// Copyright 2018 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"time"
)

// 微信绑定用户信息
type WechatUser struct {
	Id         int `xorm:"pk autoincr"`
	Openid     string
	Nickname   string
	Avatar     string
	SessionKey string
	OpenInfo   string
	Uid        int
	CreatedAt  time.Time
	UpdatedAt  time.Time `xorm:"<-"`
}
