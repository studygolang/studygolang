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
	paginator := logic.NewPaginator(goutils.MustInt(ctx.QueryParam("p"), 1))

	data := logic.DefaultIndex.FindData(ctx, tab, paginator)

	SetCookie(ctx, "INDEX_TAB", data["tab"].(string))

	data["all_nodes"] = logic.GenNodes()

	if tab == "all" {
		data["total"] = paginator.GetTotal()

	}
	return success(ctx, nil)
}
