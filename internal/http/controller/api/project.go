// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strconv"

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
	g.POST("/projects", self.Publish)
	g.GET("/projects/check_uri", self.CheckUri)
	g.GET("/projects/:uri/edit", self.Edit)
	g.PUT("/projects/:uri", self.Update)
	g.GET("/projects/:uri", self.Detail)
}

// List 项目列表
func (ProjectController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, "id DESC", "")

	total := logic.DefaultProject.Count(context.EchoContext(ctx), "")
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"projects": projects,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// Publish 发布新项目（需要登录，支持 Cookie 和 X-Token header）
func (ProjectController) Publish(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	forms, _ := ctx.FormParams()

	// 基本字段验证
	name := forms.Get("name")
	if name == "" {
		return fail(ctx, "项目名称不能为空")
	}
	src := forms.Get("src")
	if src == "" {
		return fail(ctx, "项目源码地址不能为空")
	}

	err = logic.DefaultProject.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败："+err.Error())
	}

	// 获取刚发布的项目的 URI
	uri := forms.Get("uri")
	if uri == "" {
		uri = forms.Get("name")
	}

	return success(ctx, map[string]interface{}{
		"uri": uri,
	})
}

// Edit 获取项目编辑数据（需要登录，验证权限）
func (ProjectController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	uri := ctx.Param("uri")
	project := logic.DefaultProject.FindOne(context.EchoContext(ctx), uri)
	if project == nil || project.Id == 0 {
		return fail(ctx, "项目不存在或已下线")
	}

	// 验证权限：只能编辑自己的项目，或管理员可以编辑所有项目
	if !logic.CanEdit(me, project) {
		return fail(ctx, "没有编辑权限")
	}

	return success(ctx, map[string]interface{}{
		"project": project,
	})
}

// Update 更新项目（需要登录，验证权限）
func (ProjectController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	uri := ctx.Param("uri")
	project := logic.DefaultProject.FindOne(context.EchoContext(ctx), uri)
	if project == nil || project.Id == 0 {
		return fail(ctx, "项目不存在或已下线")
	}

	// 验证权限
	if !logic.CanEdit(me, project) {
		return fail(ctx, "没有编辑权限")
	}

	forms, _ := ctx.FormParams()
	forms.Set("id", strconv.Itoa(project.Id))

	err = logic.DefaultProject.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "更新失败："+err.Error())
	}

	return success(ctx, map[string]interface{}{"uri": project.Uri})
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

// CheckUri 检查项目 URI 是否已存在
func (ProjectController) CheckUri(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	if uri == "" {
		return success(ctx, map[string]interface{}{"exists": false})
	}

	exists := logic.DefaultProject.UriExists(context.EchoContext(ctx), uri)
	return success(ctx, map[string]interface{}{"exists": exists})
}
