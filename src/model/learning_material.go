// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

type LearningMaterial struct {
	Id        int       `json:"-" xorm:"pk autoincr"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Type      int       `json:"type"`
	Seq       int       `json:"-"`
	FirstUrl  string    `json:"first_url"`
	CreatedAt time.Time `json:"-" xorm:"created"`
}
