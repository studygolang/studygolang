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

type WikiController struct{}

func (self WikiController) RegisterRoute(g *echo.Group) {
	g.GET("/wiki", self.List)
	g.GET("/wiki/:uri", self.Detail)
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

	if lastId != 0 {
		prevId = lastId
		if prevId-wikis[0].Id > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		wikis = wikis[:limit]
		nextId = wikis[limit-1].Id
	} else {
		nextId = wikis[num-1].Id
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
