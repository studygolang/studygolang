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

type WikiController struct{}

func (self WikiController) RegisterRoute(g *echo.Group) {
	g.GET("/wiki", self.List)
	g.GET("/wiki/:uri/edit", self.Edit)
	g.GET("/wiki/:uri", self.Detail)
	g.POST("/wiki", self.Create)
	g.PUT("/wiki/:id", self.Update)
}

// List Wiki 列表
func (WikiController) List(ctx echo.Context) error {
	limit := 20
	lastId := goutils.MustInt(ctx.QueryParam("lastid"))

	wikis := logic.DefaultWiki.FindBy(context.EchoContext(ctx), limit+5, lastId)
	if wikis == nil {
		return success(ctx, map[string]interface{}{
			"wikis": []*model.Wiki{},
			"page":  map[string]interface{}{},
		})
	}

	num := len(wikis)
	if num == 0 {
		return success(ctx, map[string]interface{}{
			"wikis": wikis,
			"page":  map[string]interface{}{},
		})
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	// 游标分页（id DESC）：
	// - hasNext：拉取的条数 > limit，说明还有下一页；当前页截断到 limit 条，next_id 指向当前页最后一条
	// - hasPrev：lastId > 0 表示不是第一页，prev_id 指向当前页第一条
	// 不再使用魔数 5 或 prevId += limit（在 cursor 分页中无意义）
	if num > limit {
		hasNext = true
		wikis = wikis[:limit]
		nextId = wikis[limit-1].Id
	}

	if lastId > 0 && num > 0 {
		hasPrev = true
		prevId = wikis[0].Id
	}

	return success(ctx, map[string]interface{}{
		"wikis": wikis,
		"page": map[string]interface{}{
			"has_prev": hasPrev,
			"prev_id":  prevId,
			"has_next": hasNext,
			"next_id":  nextId,
		},
	})
}

// Detail Wiki 详情（通过 URI 查询）
func (WikiController) Detail(ctx echo.Context) error {
	uri := ctx.Param("uri")
	wiki := logic.DefaultWiki.FindOne(context.EchoContext(ctx), uri)
	if wiki == nil || wiki.Id == 0 {
		return fail(ctx, "Wiki 不存在")
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		logic.Views.Incr(Request(ctx), model.TypeWiki, wiki.Id, me.Uid)
	} else {
		logic.Views.Incr(Request(ctx), model.TypeWiki, wiki.Id)
	}

	// 为了阅读数即时看到
	wiki.Viewnum++

	return success(ctx, map[string]interface{}{
		"wiki": wiki,
	})
}

// Edit 获取 Wiki 编辑数据（通过 URI 查找，需要登录和编辑权限）
func (WikiController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	uri := ctx.Param("uri")
	wiki := logic.DefaultWiki.FindOne(context.EchoContext(ctx), uri)
	if wiki == nil || wiki.Id == 0 {
		return fail(ctx, "Wiki 不存在")
	}

	if !logic.CanEdit(me, wiki) {
		return fail(ctx, "无权限编辑")
	}

	return success(ctx, map[string]interface{}{
		"wiki": map[string]interface{}{
			"id":      wiki.Id,
			"title":   wiki.Title,
			"content": wiki.Content,
			"uri":     wiki.Uri,
		},
	})
}

// Create 创建 Wiki（需要登录）
func (WikiController) Create(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	form, err := ctx.FormParams()
	if err != nil {
		return fail(ctx, "获取表单参数失败")
	}

	form.Set("uid", strconv.Itoa(me.Uid))

	err = logic.DefaultWiki.Create(context.EchoContext(ctx), me, form)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{
		"message": "创建成功",
	})
}

// Update 更新 Wiki（需要登录，通过 ID 更新）
func (WikiController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "Wiki ID 无效")
	}

	wiki := logic.DefaultWiki.FindById(context.EchoContext(ctx), id)
	if wiki == nil || wiki.Id == 0 {
		return fail(ctx, "Wiki 不存在")
	}

	if !logic.CanEdit(me, wiki) {
		return fail(ctx, "无权限编辑")
	}

	form, err := ctx.FormParams()
	if err != nil {
		return fail(ctx, "获取表单参数失败")
	}

	form.Set("id", strconv.Itoa(id))

	err = logic.DefaultWiki.Modify(context.EchoContext(ctx), me, form)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{
		"message": "更新成功",
	})
}
