// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	AnnouncementTypeNotice   = 1 // 公告
	AnnouncementTypeActivity = 2 // 活动
	AnnouncementTypeWarning  = 3 // 警告
)

// Announcement 公告
type Announcement struct {
	ID        int64     `json:"id" xorm:"pk autoincr"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      int       `json:"type"`
	Priority  int       `json:"priority"`
	StartTime time.Time `json:"start_time" xorm:"start_time"`
	EndTime   time.Time `json:"end_time" xorm:"end_time"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}
