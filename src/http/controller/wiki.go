// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"http/middleware"
	"logic"
	"model"
	"net/http"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	// logic.RegisterCommentObject(model.TypeArticle, logic.ArticleComment{})
	// logic.RegisterLikeObject(model.TypeArticle, logic.ArticleLike{})
}

type WikiController struct{}

// 注册路由
func (self WikiController) RegisterRoute(g *echo.Group) {
	g.Match([]string{"GET", "POST"}, "/wiki/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck())
	g.Match([]string{"GET", "POST"}, "/wiki/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
	g.GET("/wiki", self.ReadList)
	g.GET("/wiki/:uri", self.Detail)
}

// Create 创建wiki页
func (WikiController) Create(ctx echo.Context) error {
	title := ctx.FormValue("title")
	// 请求新建 wiki 页面
	if title == "" || ctx.Request().Method() != "POST" {
		return render(ctx, "wiki/new.html", map[string]interface{}{"activeWiki": "active"})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultWiki.Create(ctx, me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 1, "内部服务错误")
	}

	return success(ctx, nil)
}

// Modify 修改 Wiki 页
func (WikiController) Modify(ctx echo.Context) error {
	id := goutils.MustInt(ctx.FormValue("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/wiki")
	}

	if ctx.Request().Method() != "POST" {
		wiki := logic.DefaultWiki.FindById(ctx, id)
		if wiki.Id == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/wiki")
		}

		return render(ctx, "wiki/new.html", map[string]interface{}{"activeWiki": "active", "wiki": wiki})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultWiki.Modify(ctx, me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 1, "内部服务错误")
	}

	return success(ctx, nil)
}

// Detail 展示wiki页
func (WikiController) Detail(ctx echo.Context) error {
	wiki := logic.DefaultWiki.FindOne(ctx, ctx.Param("uri"))
	if wiki == nil {
		return ctx.Redirect(http.StatusSeeOther, "/wiki")
	}

	// likeFlag := 0
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		// 	likeFlag = logic.DefaultLike.HadLike(ctx, me.Uid, wiki.Id, model.TypeWiki)
		logic.Views.Incr(Request(ctx), model.TypeWiki, wiki.Id, me.Uid)
	} else {
		logic.Views.Incr(Request(ctx), model.TypeWiki, wiki.Id)
	}

	// 为了阅读数即时看到
	wiki.Viewnum++

	return render(ctx, "wiki/content.html", map[string]interface{}{"activeWiki": "active", "wiki": wiki})
}

// ReadList 获得wiki列表
func (WikiController) ReadList(ctx echo.Context) error {
	limit := 20

	lastId := goutils.MustInt(ctx.QueryParam("lastid"))
	wikis := logic.DefaultWiki.FindBy(ctx, limit+5, lastId)
	if wikis == nil {
		logger.Errorln("wiki controller: find wikis error")
		return ctx.Redirect(http.StatusSeeOther, "/wiki")
	}

	num := len(wikis)
	if num == 0 {
		if lastId == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return ctx.Redirect(http.StatusSeeOther, "/wiki")
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != 0 {
		prevId = lastId

		// 避免因为wiki下线，导致判断错误（所以 > 5）
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

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	// 获取当前用户喜欢对象信息
	// me, ok := ctx.Get("user").(*model.Me)
	// var likeFlags map[int]int
	// if ok {
	// 	likeFlags, _ = logic.DefaultLike.FindUserLikeObjects(ctx, me.Uid, model.TypeWiki, wikis[0].Id, nextId)
	// }

	return render(ctx, "wiki/list.html", map[string]interface{}{"wikis": wikis, "activeWiki": "active", "page": pageInfo})
}
