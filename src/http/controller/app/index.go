// Copyright 2018 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"logic"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"

	. "http"
)

type IndexController struct{}

// 注册路由
func (self IndexController) RegisterRoute(g *echo.Group) {
	g.GET("/home", self.Home)
	g.GET("/stat/site", self.WebsiteStat)
}

// Home 首页
func (IndexController) Home(ctx echo.Context) error {
	if len(logic.WebsiteSetting.IndexNavs) == 0 {
		return success(ctx, nil)
	}

	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = GetFromCookie(ctx, "INDEX_TAB")
	}

	if tab == "" {
		tab = logic.WebsiteSetting.IndexNavs[0].Tab
	}
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	data := logic.DefaultIndex.FindData(ctx, tab, paginator)

	SetCookie(ctx, "INDEX_TAB", data["tab"].(string))

	data["all_nodes"] = logic.GenNodes()

	if tab == "all" {
		data["total"] = paginator.GetTotal()

	}
	return success(ctx, data)
}

// WebsiteStat 网站统计信息
func (IndexController) WebsiteStat(ctx echo.Context) error {
	articleTotal := logic.DefaultArticle.Total()
	projectTotal := logic.DefaultProject.Total()
	topicTotal := logic.DefaultTopic.Total()
	cmtTotal := logic.DefaultComment.Total()
	resourceTotal := logic.DefaultResource.Total()
	bookTotal := logic.DefaultGoBook.Total()
	userTotal := logic.DefaultUser.Total()

	data := map[string]interface{}{
		"article":  articleTotal,
		"project":  projectTotal,
		"topic":    topicTotal,
		"resource": resourceTotal,
		"book":     bookTotal,
		"comment":  cmtTotal,
		"user":     userTotal,
	}

	return success(ctx, data)
}
