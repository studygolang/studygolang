// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	KeyNewUserWait = "new_user_wait" // 新用户注册多久能发布帖子，单位秒，0表示没限制
	KeyCanEditTime = "can_edit_time" // 发布后多久内能够编辑，单位秒
)

type UserSetting struct {
	Id        int `xorm:"pk autoincr"`
	Key       string
	Value     int
	CreatedAt time.Time `xorm:"created"`
}
