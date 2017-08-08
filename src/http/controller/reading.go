// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

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
	g.GET("/readings", self.ReadingList)
	g.GET("/readings/:id", self.IReading)
}

// ReadingList 晨读列表页
func (ReadingController) ReadingList(ctx echo.Context) error {
	limit := 20
	lastId := goutils.MustInt(ctx.QueryParam("lastid"))
	rtype := goutils.MustInt(ctx.QueryParam("rtype"), model.RtypeGo)

	readings := logic.DefaultReading.FindBy(ctx, limit+5, rtype, lastId)
	num := len(readings)
	if num == 0 {
		if lastId == 0 {
			return render(ctx, "readings/list.html", map[string]interface{}{"activeReadings": "active", "readings": readings, "rtype": rtype})
		} else {
			return ctx.Redirect(http.StatusSeeOther, "/readings")
		}
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId > 0 {
		prevId = lastId

		// 避免因为项目下线，导致判断错误（所以 > 5）
		if prevId-readings[0].Id > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		readings = readings[:limit]
		nextId = readings[limit-1].Id
	} else {
		nextId = readings[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}
	return render(ctx, "readings/list.html", map[string]interface{}{"activeReadings": "active", "readings": readings, "page": pageInfo, "rtype": rtype})
}

// IReading 点击 【我要晨读】，记录点击数，跳转
func (ReadingController) IReading(ctx echo.Context) error {
	uri := logic.DefaultReading.IReading(ctx, goutils.MustInt(ctx.Param("id")))
	return ctx.Redirect(http.StatusSeeOther, uri)
}
