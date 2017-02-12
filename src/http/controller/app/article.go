// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"logic"
	"model"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type ArticleController struct{}

// 注册路由
func (this *ArticleController) RegisterRoute(g *echo.Group) {
	g.Get("/articles", this.ReadList)
	g.Get("/article/detail", this.Detail)
}

// ReadList 网友文章列表页
func (ArticleController) ReadList(ctx echo.Context) error {
	limit := 20

	lastId := goutils.MustInt(ctx.QueryParam("base"))
	articles := logic.DefaultArticle.FindBy(ctx, limit+5, lastId)
	if articles == nil {
		return fail(ctx, "获取失败")
	}

	hasMore := false
	if len(articles) > limit {
		hasMore = true
		articles = articles[:limit]
	}

	data := map[string]interface{}{
		"has_more": hasMore,
	}

	articleList := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		articleList[i] = map[string]interface{}{
			"id":       article.Id,
			"name":     article.Name,
			"title":    article.Title,
			"pub_date": article.PubDate,
			"tags":     article.Tags,
			"viewnum":  article.Viewnum,
			"cmtnum":   article.Cmtnum,
			"likenum":  article.Likenum,
			"top":      article.Top,
			"author":   article.AuthorTxt,
		}
	}
	data["articles"] = articleList

	return success(ctx, data)
}

// Detail 文章详细页
func (ArticleController) Detail(ctx echo.Context) error {
	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(ctx, goutils.MustInt(ctx.QueryParam("id")))
	if err != nil {
		return fail(ctx, err.Error())
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return success(ctx, map[string]interface{}{"article": map[string]interface{}{"id": 0}})
	}

	logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)

	// 为了阅读数即时看到
	article.Viewnum++

	article.Txt = ""
	data := map[string]interface{}{
		"article": article,
	}

	// TODO: 暂时不用
	_ = prevNext
	return success(ctx, data)
}
