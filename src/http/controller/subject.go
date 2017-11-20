// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"http/middleware"
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type SubjectController struct{}

// 注册路由
func (self SubjectController) RegisterRoute(g *echo.Group) {
	g.Get("/subject/:id", self.Index)
	g.Post("/subject/follow", self.Follow, middleware.NeedLogin())
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

	data := map[string]interface{}{
		"subject":      subject,
		"articles":     articles,
		"article_num":  articleNum,
		"followers":    followers,
		"follower_num": followerNum,
		"order_by":     orderBy,
	}

	return render(ctx, "subject/index.html", data)
}

func (self SubjectController) Follow(ctx echo.Context) error {

}
