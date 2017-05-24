// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

type Advertisement struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	AdType    int       `json:"ad_type"`
	Code      string    `json:"code"`
	Source    string    `json:"source"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}

type PageAd struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Path      string    `json:"path"`
	AdId      int       `json:"ad_id"`
	Position  string    `json:"position"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at" xorm:"<-"`
}
