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

	// 支持排序参数：hot（热门 star+watch）、latest（最新 id）、noreply（零回复 cmtnum）
	sort := ctx.QueryParam("sort")
	var orderBy string
	switch sort {
	case "hot":
		// OpenProject 没有 star/watch/fork 列，用 likenum（点赞）+ viewnum（浏览）表示热度
		orderBy = "likenum DESC, viewnum DESC, id DESC"
	case "noreply":
		orderBy = "cmtnum ASC, id DESC"
	default: // latest 或空值
		orderBy = "id DESC"
	}

	// 只展示在线和新建状态的项目（过滤已下线的）
	projectFilter := "status IN(?,?)"
	projectFilterArgs := []interface{}{model.ProjectStatusNew, model.ProjectStatusOnline}

	projects := logic.DefaultProject.FindAll(context.EchoContext(ctx), paginator, orderBy, projectFilter, projectFilterArgs...)

	total := logic.DefaultProject.Count(context.EchoContext(ctx), projectFilter, projectFilterArgs...)
	hasMore := paginator.SetTotal(total).HasMorePage()

	// Email 脱敏：OpenProject.User / LastReplyUser 字段会自动序列化暴露 Email
	sanitizeProjectsForPublic(projects)

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

	// 敏感词检查
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}
	// 余额检查
	if !balanceCheck(me, false) {
		return failBalance(ctx)
	}
	// 验证码检查
	if !captchaCheck(ctx, me) {
		return failCaptcha(ctx)
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

	// 发布后邮件通知站长
	publishNotice(ctx, me)

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

	// Email 脱敏：OpenProject.User / LastReplyUser 字段会自动序列化暴露 Email
	sanitizeProjectForPublic(project)

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

	// 敏感词检查（与 master /projects/modify 中间件保持一致；
	// 漏挂会导致更新时塞入广告内容不被检测/冻结）
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
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

	// Email 脱敏：OpenProject.User / LastReplyUser 字段会自动序列化暴露 Email
	sanitizeProjectForPublic(project)

	result := map[string]interface{}{
		"project": project,
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		result["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, project.Id, model.TypeProject)
		result["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, project.Id, model.TypeProject)

		logic.Views.Incr(Request(ctx), model.TypeProject, project.Id, me.Uid)

		if me.Uid != project.User.Uid {
			go logic.DefaultViewRecord.Record(project.Id, model.TypeProject, me.Uid)
		}

		if me.IsRoot || me.Uid == project.User.Uid {
			result["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), project.Id, model.TypeProject)
			result["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), project.Id, model.TypeProject)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeProject, project.Id)
	}

	// 为了阅读数即时看到
	project.Viewnum++

	// 评论信息
	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), project.Id, model.TypeProject, 0, project.Lastreplyuid,
	)
	if project.Lastreplyuid != 0 {
		// Email 脱敏：未公开邮箱的用户清空 Email，与 profile.html 语义对齐
		project.LastReplyUser = sanitizeUserForPublic(lastReplyUser)
	}
	result["replies"] = normalizeReplies(replies)

	return success(ctx, result)
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
