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

type ResourceController struct{}

func (self ResourceController) RegisterRoute(g *echo.Group) {
	g.GET("/resources/categories", self.Categories)
	g.GET("/resources", self.List)
	g.POST("/resources", self.Publish)
	g.GET("/resources/:id/edit", self.Edit)
	g.PUT("/resources/:id", self.Update)
	g.GET("/resources/:id", self.Detail)
}

// List 资源列表（支持 catid 分类过滤和 sort 排序）
// sort 可选值：new（最新，默认）、hot（最热：likenum + viewnum）、recommend（推荐：cmtnum）
func (ResourceController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	catid := goutils.MustInt(ctx.QueryParam("catid"))
	sort := ctx.QueryParam("sort")

	order := resolveResourceOrder(sort)

	var (
		resources interface{}
		total     int64
	)

	if catid > 0 {
		resources, total = logic.DefaultResource.FindByCatid(context.EchoContext(ctx), paginator, catid, order)
	} else {
		resources, total = logic.DefaultResource.FindAll(context.EchoContext(ctx), paginator, order, "")
	}

	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"resources": resources,
		"has_more":  hasMore,
		"total":     total,
		"page":      curPage,
		"sort":      sort,
	})
}

// resolveResourceOrder 将 sort 参数转换为数据库排序表达式
func resolveResourceOrder(sort string) string {
	switch sort {
	case "hot":
		// 按点赞数 + 浏览数综合排序
		return "resource_ex.likenum DESC, resource_ex.viewnum DESC"
	case "recommend":
		// 按评论数排序
		return "resource_ex.cmtnum DESC"
	default:
		// 默认按更新时间倒序
		return "resource.mtime DESC"
	}
}

// Publish 发布新资源（需要登录，支持 Cookie 和 X-Token header）
func (ResourceController) Publish(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	forms, _ := ctx.FormParams()

	// 基本字段验证
	title := forms.Get("title")
	if title == "" {
		return fail(ctx, "资源标题不能为空")
	}
	form := forms.Get("form")
	if form == "" {
		return fail(ctx, "请选择资源类型")
	}
	// 验证 form 值合法
	if form != model.LinkForm && form != model.ContentForm {
		return fail(ctx, "资源类型不合法")
	}
	// 链接类型需要 url，内容类型需要 content
	if form == model.LinkForm && forms.Get("url") == "" {
		return fail(ctx, "链接类型资源需要填写 URL")
	}
	if form == model.ContentForm && forms.Get("content") == "" {
		return fail(ctx, "内容类型资源需要填写内容")
	}

	err = logic.DefaultResource.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败："+err.Error())
	}

	// 获取刚发布的资源 ID（从逻辑层返回的 form 中获取）
	id := forms.Get("id")

	return success(ctx, map[string]interface{}{
		"id": id,
	})
}

// Edit 获取资源编辑数据（需要登录，验证权限）
func (ResourceController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "资源 ID 非法")
	}

	resource := logic.DefaultResource.FindResource(context.EchoContext(ctx), id)
	if resource == nil || resource.Id == 0 {
		return fail(ctx, "资源不存在")
	}

	// 验证权限
	if !logic.CanEdit(me, resource) {
		return fail(ctx, "没有编辑权限")
	}

	return success(ctx, map[string]interface{}{
		"resource": resource,
	})
}

// Update 更新资源（需要登录，验证权限）
func (ResourceController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "资源 ID 非法")
	}

	resource := logic.DefaultResource.FindResource(context.EchoContext(ctx), id)
	if resource == nil || resource.Id == 0 {
		return fail(ctx, "资源不存在")
	}

	// 验证权限
	if !logic.CanEdit(me, resource) {
		return fail(ctx, "没有编辑权限")
	}

	forms, _ := ctx.FormParams()
	forms.Set("id", strconv.Itoa(resource.Id))

	err = logic.DefaultResource.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "更新失败："+err.Error())
	}

	return success(ctx, map[string]interface{}{"id": resource.Id})
}

// Detail 资源详情
func (ResourceController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "参数有误")
	}

	resource, comments := logic.DefaultResource.FindById(context.EchoContext(ctx), id)
	if len(resource) == 0 {
		return fail(ctx, "获取失败")
	}

	logic.Views.Incr(Request(ctx), model.TypeResource, id)

	return success(ctx, map[string]interface{}{
		"resource": resource,
		"comments": comments,
	})
}

// Categories 获取资源分类列表
func (ResourceController) Categories(ctx echo.Context) error {
	return success(ctx, map[string]interface{}{
		"categories": logic.AllCategory,
	})
}
