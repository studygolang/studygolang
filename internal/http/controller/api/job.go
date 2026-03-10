// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type JobController struct{}

func (self JobController) RegisterRoute(g *echo.Group) {
	g.GET("/jobs", self.List)
}

// List 职位列表
func (JobController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}

	jobs, total, err := logic.DefaultJob.FindAll(context.EchoContext(ctx), curPage, perPage)
	if err != nil {
		// 数据库表不存在时返回空列表，不报错
		return success(ctx, map[string]interface{}{
			"jobs":     make([]interface{}, 0),
			"total":    0,
			"page":     curPage,
			"has_more": false,
		})
	}

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"jobs":     jobs,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}
