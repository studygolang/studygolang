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

type ArticleController struct{}

// RegisterRoute 注册路由
func (self *ArticleController) RegisterRoute(g *echo.Group) {
	g.GET("/articles", self.List)
	g.GET("/articles/:id", self.Detail)
}

// List 文章列表，支持 p 分页参数
func (ArticleController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	topArticles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "top=1")
	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "")

	total := logic.DefaultArticle.Count(context.EchoContext(ctx), "")
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"list":     append(topArticles, articles...),
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// Detail 文章详情，增加浏览量
func (ArticleController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "文章 id 非法")
	}

	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(context.EchoContext(ctx), id)
	if err != nil {
		return fail(ctx, err.Error())
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return success(ctx, map[string]interface{}{"article": map[string]interface{}{"id": 0}})
	}

	logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)
	article.Viewnum++

	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), article.Id, model.TypeArticle, 0, article.Lastreplyuid,
	)
	if article.Lastreplyuid != 0 {
		article.LastReplyUser = lastReplyUser
	}

	article.Txt = ""

	return success(ctx, map[string]interface{}{
		"article":   article,
		"replies":   replies,
		"prev_next": prevNext,
	})
}
