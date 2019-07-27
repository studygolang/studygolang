// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

type ViewRecord struct {
	Id        int       `json:"id" xorm:"pk autoincr"`
	Objid     int       `json:"objid"`
	Objtype   int       `json:"objtype"`
	Uid       int       `json:"uid"`
	CreatedAt OftenTime `json:"created_at" xorm:"<-"`
}
