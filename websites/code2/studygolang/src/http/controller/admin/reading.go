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
	"github.com/polaris1119/goutils"
)

type ReadingController struct{}

// 注册路由
func (self ReadingController) RegisterRoute(g *echo.Group) {
	g.GET("/reading/list", self.ReadingList)
	g.POST("/reading/query.html", self.ReadingQuery)
	g.Match([]string{"GET", "POST"}, "/reading/publish", self.Publish)
}

// ReadingList 所有晨读（分页）
func (ReadingController) ReadingList(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)

	readings, total := logic.DefaultReading.FindReadingByPage(ctx, nil, curPage, limit)
	if readings == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   readings,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "reading/list.html,reading/query.html", data)
}

// ReadingQuery
func (ReadingController) ReadingQuery(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	conds := parseConds(ctx, []string{"id", "rtype"})

	readings, total := logic.DefaultReading.FindReadingByPage(ctx, conds, curPage, limit)
	if readings == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   readings,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "reading/query.html", data)
}

// Publish
func (ReadingController) Publish(ctx echo.Context) error {
	var data = make(map[string]interface{})

	if ctx.FormValue("submit") == "1" {
		user := ctx.Get("user").(*model.Me)
		errMsg, err := logic.DefaultReading.SaveReading(ctx, ctx.FormParams(), user.Username)
		if err != nil {
			return fail(ctx, 1, errMsg)
		}
		return success(ctx, nil)
	}

	id := goutils.MustInt(ctx.QueryParam("id"))
	if id != 0 {
		reading := logic.DefaultReading.FindById(ctx, id)
		if reading != nil {
			data["reading"] = reading
		}
	}

	return render(ctx, "reading/modify.html", data)
}
