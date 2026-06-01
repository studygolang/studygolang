// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"context"
	"time"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/internal/model"

	"github.com/polaris1119/logger"
)

type AnnouncementLogic struct{}

var DefaultAnnouncement = AnnouncementLogic{}

// FindById 根据 id 查找公告
func (self AnnouncementLogic) FindById(ctx context.Context, id int64) (*model.Announcement, error) {
	announcement := &model.Announcement{ID: id}
	has, err := MasterDB.Get(announcement)
	if err != nil {
		logger.Errorln("AnnouncementLogic FindById error:", err)
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return announcement, nil
}

// FindActive 根据 type 过滤，返回当前时间在 [start_time, end_time] 区间的公告列表，
// 按 priority DESC, created_at DESC 排序。
func (self AnnouncementLogic) FindActive(ctx context.Context, annType int, paginator *Paginator) ([]*model.Announcement, int64) {
	now := time.Now()

	session := MasterDB.Where("start_time <= ? AND end_time >= ?", now, now)
	if annType > 0 {
		session = session.And("type = ?", annType)
	}

	total, err := session.Count(new(model.Announcement))
	if err != nil {
		logger.Errorln("AnnouncementLogic FindActive count error:", err)
		return nil, 0
	}

	var announcements []*model.Announcement
	err = session.OrderBy("priority DESC, created_at DESC").
		Limit(paginator.PerPage(), paginator.Offset()).
		Find(&announcements)
	if err != nil {
		logger.Errorln("AnnouncementLogic FindActive find error:", err)
		return nil, 0
	}

	return announcements, total
}
