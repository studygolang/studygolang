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

type AnnouncementController struct{}

// RegisterRoute 注册路由
func (self AnnouncementController) RegisterRoute(g *echo.Group) {
	g.GET("/announcements", self.List)
	g.GET("/announcements/:id", self.Detail)
}

// List 获取公告列表
// 支持参数:
//   - p: 页码，默认 1
//   - type: 公告类型筛选（1=公告，2=活动，3=警告）
//
// 返回当前时间在 [start_time, end_time] 区间的公告，按 priority DESC, created_at DESC 排序
func (AnnouncementController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	annType := goutils.MustInt(ctx.QueryParam("type"), 0)

	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	announcements, total := logic.DefaultAnnouncement.FindActive(
		context.EchoContext(ctx), annType, paginator,
	)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"list":     announcements,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// Detail 获取公告详情
func (AnnouncementController) Detail(ctx echo.Context) error {
	id := goutils.MustInt64(ctx.Param("id"), 0)
	if id <= 0 {
		return fail(ctx, "无效的公告ID")
	}

	announcement, err := logic.DefaultAnnouncement.FindById(context.EchoContext(ctx), id)
	if err != nil {
		return fail(ctx, "获取公告失败")
	}
	if announcement == nil {
		return fail(ctx, "公告不存在", 404)
	}

	return success(ctx, map[string]interface{}{
		"announcement": announcement,
	})
}
