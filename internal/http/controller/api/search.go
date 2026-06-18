// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/http"
	"net/url"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type SearchController struct{}

func (self SearchController) RegisterRoute(g *echo.Group) {
	g.GET("/search", self.Search)
	g.GET("/tag/:name", self.TagList)
}

// Search 全文搜索（q 关键词，p 分页，type 内容类型）
func (SearchController) Search(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	if q == "" {
		return fail(ctx, "搜索关键词不能为空")
	}
	// 限制关键词长度，避免超长 query 打满搜索引擎（master 也未限制，refactor 顺手补上）
	if len([]rune(q)) > 64 {
		return fail(ctx, "搜索关键词过长")
	}

	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	field := ctx.QueryParam("type")

	rows := 50
	respBody, err := logic.DefaultSearcher.DoSearch(q, field, (p-1)*rows, rows)
	if err != nil {
		return fail(ctx, "搜索服务异常")
	}

	paginator := logic.NewPaginatorWithPerPage(p, rows)
	hasMore := paginator.SetTotal(int64(respBody.NumFound)).HasMorePage()

	// 将 Docs 转换为前端期望的 results 数组，字段名与 SearchResult 类型对齐
	results := respBody.Docs
	if results == nil {
		results = make([]*model.Document, 0)
	}

	return success(ctx, map[string]interface{}{
		"results":  results,
		"keyword":  q,
		"type":     field,
		"page":     p,
		"has_more": hasMore,
		"total":    respBody.NumFound,
	})
}

// TagList 标签内容列表（复用搜索功能，field=tag）
func (SearchController) TagList(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	var err error
	name, err = url.QueryUnescape(name)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	// 限制 tag 长度（防止滥用）
	if len([]rune(name)) > 9 {
		return fail(ctx, "标签名称过长")
	}

	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	rows := 50

	respBody, err := logic.DefaultSearcher.DoSearch(name, "tag", (p-1)*rows, rows)
	if err != nil {
		return fail(ctx, "搜索服务异常")
	}

	_, nodes := logic.DefaultSearcher.FillNodeAndUser(context.EchoContext(ctx), respBody)

	paginator := logic.NewPaginatorWithPerPage(p, rows)
	hasMore := paginator.SetTotal(int64(respBody.NumFound)).HasMorePage()

	results := respBody.Docs
	if results == nil {
		results = make([]*model.Document, 0)
	}

	return success(ctx, map[string]interface{}{
		"results":  results,
		"keyword":  name,
		"nodes":    nodes,
		"page":     p,
		"has_more": hasMore,
		"total":    respBody.NumFound,
	})
}
