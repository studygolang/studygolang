// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type ProjectController struct{}

func (self ProjectController) RegisterRoute(g *echo.Group) {
	g.GET("/projects", self.List)
	g.GET("/projects/:uri", self.Detail)
}

// List 项目列表
func (ProjectController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, "id DESC", "")

	total := logic.DefaultProject.Count(context.EchoContext(ctx), "")
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"projects": projects,
		"has_more": hasMore,
	})
}

// Detail 项目详情（通过 URI 查询）
func (ProjectController) Detail(ctx echo.Context) error {
	uri := ctx.Param("uri")
	project := logic.DefaultProject.FindOne(context.EchoContext(ctx), uri)
	if project == nil || project.Id == 0 {
		return fail(ctx, "获取失败或已下线")
	}

	logic.Views.Incr(Request(ctx), model.TypeProject, project.Id)

	// 为了阅读数即时看到
	project.Viewnum++

	// 评论信息
	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), project.Id, model.TypeProject, 0, project.Lastreplyuid,
	)
	if project.Lastreplyuid != 0 {
		project.LastReplyUser = lastReplyUser
	}

	return success(ctx, map[string]interface{}{
		"project": project,
		"replies": replies,
	})
}
