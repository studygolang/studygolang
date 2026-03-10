// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"context"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/internal/model"
)

type JobLogic struct{}

var DefaultJob = JobLogic{}

// FindAll 获取职位列表
func (JobLogic) FindAll(ctx context.Context, curPage, pageSize int) ([]*model.Job, int64, error) {
	objLog := GetLogger(ctx)

	jobs := make([]*model.Job, 0)
	total, err := MasterDB.Where("status=0").
		Desc("id").
		Limit(pageSize, (curPage-1)*pageSize).
		FindAndCount(&jobs)
	if err != nil {
		objLog.Errorln("JobLogic FindAll error:", err)
		return nil, 0, err
	}
	return jobs, total, nil
}
