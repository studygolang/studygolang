// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"

	. "http"
	"model"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeArticle, logic.ArticleComment{})
	logic.RegisterLikeObject(model.TypeArticle, logic.ArticleLike{})
}

type ArticleController struct{}

// 注册路由
func (this *ArticleController) RegisterRoute(g *echo.Group) {
	g.Get("/articles", this.ReadList)
	g.Get("/articles/:id", this.Detail)
}

// ReadList 网友文章列表页
func (ArticleController) ReadList(ctx echo.Context) error {
	limit := 20

	lastId := goutils.MustInt(ctx.QueryParam("lastid"))
	articles := logic.DefaultArticle.FindBy(ctx, limit+5, lastId)
	if articles == nil {
		logger.Errorln("article controller: find article error")
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	num := len(articles)
	if num == 0 {
		if lastId == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != 0 {
		prevId = lastId

		firstNoTopId := articles[0].Id
		for i := 0; i < num; i++ {
			if articles[i].Top != 1 {
				firstNoTopId = articles[i].Id
				break
			}
		}
		// 避免因为文章下线，导致判断错误（所以 > 5）
		if prevId-firstNoTopId > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		articles = articles[:limit]
		nextId = articles[limit-1].Id
	} else {
		nextId = articles[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	// 获取当前用户喜欢对象信息
	me, ok := ctx.Get("user").(*model.Me)
	var likeFlags map[int]int
	if ok {
		likeFlags, _ = logic.DefaultLike.FindUserLikeObjects(ctx, me.Uid, model.TypeArticle, articles[0].Id, nextId)
	}

	return render(ctx, "articles/list.html", map[string]interface{}{"articles": articles, "activeArticles": "active", "page": pageInfo, "likeflags": likeFlags})
}

// Detail 文章详细页
func (ArticleController) Detail(ctx echo.Context) error {
	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(ctx, goutils.MustInt(ctx.Param("id")))
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	likeFlag := 0
	hadCollect := 0
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		likeFlag = logic.DefaultLike.HadLike(ctx, me.Uid, article.Id, model.TypeArticle)
		hadCollect = logic.DefaultFavorite.HadFavorite(ctx, me.Uid, article.Id, model.TypeArticle)
	}

	logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)

	// 为了阅读数即时看到
	article.Viewnum++

	return render(ctx, "articles/detail.html,common/comment.html", map[string]interface{}{"activeArticles": "active", "article": article, "prev": prevNext[0], "next": prevNext[1], "likeflag": likeFlag, "hadcollect": hadCollect})
}
