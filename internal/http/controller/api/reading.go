// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type ReadingController struct{}

func (self ReadingController) RegisterRoute(g *echo.Group) {
	g.GET("/readings", self.List)
	g.GET("/readings/:id", self.Detail)
}

// List 晨读列表（支持 lastid、rtype 参数）
func (ReadingController) List(ctx echo.Context) error {
	limit := 20
	lastId := goutils.MustInt(ctx.QueryParam("lastid"))
	rtype := goutils.MustInt(ctx.QueryParam("rtype"), model.RtypeGo)

	readings := logic.DefaultReading.FindBy(context.EchoContext(ctx), limit+5, rtype, lastId)
	num := len(readings)

	if num == 0 {
		return success(ctx, map[string]interface{}{
			"readings": readings,
			"rtype":    rtype,
			"page":     map[string]interface{}{},
		})
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId > 0 {
		prevId = lastId
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

	return success(ctx, map[string]interface{}{
		"readings": readings,
		"rtype":    rtype,
		"page": map[string]interface{}{
			"has_prev": hasPrev,
			"prev_id":  prevId,
			"has_next": hasNext,
			"next_id":  nextId,
		},
	})
}

// Detail 晨读详情
func (ReadingController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "参数有误")
	}

	reading := logic.DefaultReading.FindById(context.EchoContext(ctx), id)
	if reading == nil || reading.Id == 0 {
		return fail(ctx, "晨读不存在")
	}

	return success(ctx, map[string]interface{}{
		"reading": reading,
	})
}
