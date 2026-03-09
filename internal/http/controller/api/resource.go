// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type ResourceController struct{}

func (self ResourceController) RegisterRoute(g *echo.Group) {
	g.GET("/resources", self.List)
	g.GET("/resources/:id", self.Detail)
}

// List 资源列表（支持 catid 分类过滤）
func (ResourceController) List(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	catid := goutils.MustInt(ctx.QueryParam("catid"))

	var (
		resources interface{}
		total     int64
	)

	if catid > 0 {
		resources, total = logic.DefaultResource.FindByCatid(context.EchoContext(ctx), paginator, catid)
	} else {
		resources, total = logic.DefaultResource.FindAll(context.EchoContext(ctx), paginator, "resource.mtime", "")
	}

	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"resources": resources,
		"has_more":  hasMore,
	})
}

// Detail 资源详情
func (ResourceController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.Param("id"))
	if id == 0 {
		return fail(ctx, "参数有误")
	}

	resource, comments := logic.DefaultResource.FindById(context.EchoContext(ctx), id)
	if len(resource) == 0 {
		return fail(ctx, "获取失败")
	}

	logic.Views.Incr(Request(ctx), model.TypeResource, id)

	return success(ctx, map[string]interface{}{
		"resource": resource,
		"comments": comments,
	})
}
