// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

type DefaultAvatar struct {
	Id        int `json:"-" xorm:"pk autoincr"`
	Filename  string
	CreatedAt time.Time `json:"-" xorm:"<-"`
}
