// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"sort"
	"strconv"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type WikiController struct{}

// publicWiki 把 logic 填充的 Users（完整 model.User，含 Email/Balance/IsRoot 等）
// 投影为公开字段：author（所有者用户名）+ contributors（uid/username/avatar）。
// Users 已带 json:"-" 不会序列化，这里再置空避免内存中残留敏感数据被其他编码路径输出。
func publicWiki(wiki *model.Wiki) {
	if wiki == nil {
		return
	}
	if owner, ok := wiki.Users[wiki.Uid]; ok && owner != nil {
		wiki.Author = owner.Username
	}
	contributors := make([]model.WikiContributor, 0, len(wiki.Users))
	for _, u := range wiki.Users {
		if u == nil {
			continue
		}
		contributors = append(contributors, model.WikiContributor{
			Uid:      u.Uid,
			Username: u.Username,
			Avatar:   u.Avatar,
		})
	}
	sort.Slice(contributors, func(i, j int) bool { return contributors[i].Uid < contributors[j].Uid })
	wiki.Contributors = contributors
	wiki.Users = nil
}

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
	previous := lastId > 0 && ctx.QueryParam("direction") == "prev"

	var wikis []*model.Wiki
	if previous {
		wikis = logic.DefaultWiki.FindPrev(context.EchoContext(ctx), limit+1, lastId)
	} else {
		wikis = logic.DefaultWiki.FindBy(context.EchoContext(ctx), limit+1, lastId)
	}
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

	hasPrev := lastId > 0 && !previous
	hasNext := previous
	if num > limit {
		if previous {
			hasPrev = true
			wikis = wikis[num-limit:]
		} else {
			hasNext = true
			wikis = wikis[:limit]
		}
	}

	prevId := wikis[0].Id
	nextId := wikis[len(wikis)-1].Id
	for _, wiki := range wikis {
		publicWiki(wiki)
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
	uri := PathParam(ctx, "uri")
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

	publicWiki(wiki)

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

	uri := PathParam(ctx, "uri")
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
// master 在路由层挂了 Sensivite + BalanceCheck，refactor 漏挂；此处补回。
func (WikiController) Create(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	// 敏感词检查（与 master /wiki/new 中间件保持一致）
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}
	// 余额检查（与 master /wiki/new 中间件保持一致）
	if !balanceCheck(me, false) {
		return failBalance(ctx)
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
// master 在路由层挂了 Sensivite，refactor 漏挂；此处补回。
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

	// 敏感词检查（与 master /wiki/modify 中间件保持一致）
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
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
