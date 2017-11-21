// Copyright 2017 The StudyGolang Authors. All rights reserved.
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

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type SubjectController struct{}

// 注册路由
func (self SubjectController) RegisterRoute(g *echo.Group) {
	g.Get("/subject/:id", self.Index)
	g.Post("/subject/follow", self.Follow, middleware.NeedLogin())
	g.Get("/subject/my_articles", self.MyArticles, middleware.NeedLogin())
	g.Post("/subject/contribute", self.Contribute, middleware.NeedLogin())
	g.Post("/subject/remove_contribute", self.RemoveContribute, middleware.NeedLogin())
}

func (SubjectController) Index(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	subject := logic.DefaultSubject.FindOne(ctx, id)
	if subject.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	orderBy := ctx.QueryParam("order_by")
	articles := logic.DefaultSubject.FindArticles(ctx, id, orderBy)
	if orderBy == "" {
		orderBy = "added_at"
	}

	articleNum := logic.DefaultSubject.FindArticleTotal(ctx, id)

	followers := logic.DefaultSubject.FindFollowers(ctx, id)
	followerNum := logic.DefaultSubject.FindFollowerTotal(ctx, id)

	// 是否已关注
	followed := false
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		followed = logic.DefaultSubject.HadFollow(ctx, id, me)
	}

	data := map[string]interface{}{
		"subject":      subject,
		"articles":     articles,
		"article_num":  articleNum,
		"followers":    followers,
		"follower_num": followerNum,
		"order_by":     orderBy,
		"followed":     followed,
	}

	return render(ctx, "subject/index.html", data)
}

func (self SubjectController) Follow(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultSubject.Follow(ctx, sid, me)
	if err != nil {
		return fail(ctx, 1, "关注失败！")
	}

	return success(ctx, nil)
}

func (self SubjectController) MyArticles(ctx echo.Context) error {
	kw := ctx.QueryParam("kw")
	sid := goutils.MustInt(ctx.FormValue("sid"))

	me := ctx.Get("user").(*model.Me)

	articles := logic.DefaultArticle.SearchMyArticles(ctx, me, sid, kw)

	return success(ctx, map[string]interface{}{
		"articles": articles,
	})
}

// Contribute 投稿
func (self SubjectController) Contribute(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	me := ctx.Get("user").(*model.Me)

	err := logic.DefaultSubject.Contribute(ctx, me, sid, articleId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

// RemoveContribute 删除投稿
func (self SubjectController) RemoveContribute(ctx echo.Context) error {
	sid := goutils.MustInt(ctx.FormValue("sid"))
	articleId := goutils.MustInt(ctx.FormValue("article_id"))

	err := logic.DefaultSubject.RemoveContribute(ctx, sid, articleId)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}
