// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

type FriendLink struct {
	Id        int       `json:"-" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	Logo      string    `json:"logo"`
	Seq       int       `json:"-"`
	CreatedAt time.Time `json:"-" xorm:"created"`
}
