// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/http"
	"github.com/studygolang/studygolang/logic"
	"github.com/studygolang/studygolang/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type ArticleController struct{}

// 注册路由
func (this *ArticleController) RegisterRoute(g *echo.Group) {
	g.GET("/articles", this.ReadList)
	g.GET("/article/detail", this.Detail)
}

// ReadList 网友文章列表页
func (ArticleController) ReadList(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	// 置顶的 article
	topArticles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "top=1")

	articles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "")

	total := logic.DefaultArticle.Count(context.EchoContext(ctx), "")
	hasMore := paginator.SetTotal(total).HasMorePage()

	data := map[string]interface{}{
		"articles": append(topArticles, articles...),
		"has_more": hasMore,
	}

	return success(ctx, data)
}

// Detail 文章详细页
func (ArticleController) Detail(ctx echo.Context) error {
	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(context.EchoContext(ctx), goutils.MustInt(ctx.QueryParam("id")))
	if err != nil {
		return fail(ctx, err.Error())
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return success(ctx, map[string]interface{}{"article": map[string]interface{}{"id": 0}})
	}

	logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)

	// 为了阅读数即时看到
	article.Viewnum++

	// 回复信息（评论）
	replies, _, lastReplyUser := logic.DefaultComment.FindObjComments(context.EchoContext(ctx), article.Id, model.TypeArticle, 0, article.Lastreplyuid)
	// 有人回复
	if article.Lastreplyuid != 0 {
		article.LastReplyUser = lastReplyUser
	}

	article.Txt = ""
	data := map[string]interface{}{
		"article": article,
		"replies": replies,
	}

	// TODO: 暂时不用
	_ = prevNext
	return success(ctx, data)
}
