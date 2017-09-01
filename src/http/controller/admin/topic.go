// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"logic"
	"model"
	"net/http"

	"github.com/labstack/echo"
)

type TopicController struct{}

// 注册路由
func (self TopicController) RegisterRoute(g *echo.Group) {
	g.GET("/community/topic/list", self.List)
	g.POST("/community/topic/query.html", self.Query)
	g.Match([]string{"GET", "POST"}, "/community/topic/modify", self.Modify)
}

// List 所有主题（分页）
func (TopicController) List(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	topics, total := logic.DefaultTopic.FindByPage(ctx, nil, curPage, limit)

	if topics == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   topics,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "topic/list.html,topic/query.html", data)
}

// Query
func (TopicController) Query(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	conds := parseConds(ctx, []string{"tid", "title", "uid"})

	articles, total := logic.DefaultTopic.FindByPage(ctx, conds, curPage, limit)

	if articles == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   articles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return renderQuery(ctx, "topic/query.html", data)
}

// Modify
func (self TopicController) Modify(ctx echo.Context) error {
	var data = make(map[string]interface{})

	if ctx.FormValue("submit") == "1" {
		user := ctx.Get("user").(*model.Me)
		errMsg, err := logic.DefaultArticle.Modify(ctx, user, ctx.FormParams())
		if err != nil {
			return fail(ctx, 1, errMsg)
		}
		return success(ctx, nil)
	}
	article, err := logic.DefaultArticle.FindById(ctx, ctx.QueryParam("id"))
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, ctx.Echo().URI(echo.HandlerFunc(self.List)))
	}

	data["article"] = article
	data["statusSlice"] = model.ArticleStatusSlice
	data["langSlice"] = model.LangSlice

	return render(ctx, "topic/modify.html", data)
}
