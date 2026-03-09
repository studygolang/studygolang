// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type IndexController struct{}

// RegisterRoute 注册路由
func (self IndexController) RegisterRoute(g *echo.Group) {
	g.GET("/home", self.Home)
	g.GET("/stat/site", self.WebsiteStat)
}

// Home 首页话题列表，支持 tab 参数切换分类
// 返回格式与 TopicListData 对齐：topics/tab/tab_list/total/page/has_more
func (IndexController) Home(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = "all"
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	var (
		topTopics []map[string]interface{}
		topics    []map[string]interface{}
		total     int64
	)

	hotNodes := logic.DefaultTopic.FindHotNodes(context.EchoContext(ctx))

	switch tab {
	case "all":
		topTopics = logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "ctime DESC", "top=1")
		topics = logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", "top!=1")
		total = logic.DefaultTopic.Count(context.EchoContext(ctx), "top!=1")
	case "no_reply":
		topics = logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", "lastreplyuid=?", 0)
		total = logic.DefaultTopic.Count(context.EchoContext(ctx), "lastreplyuid=?", 0)
	default:
		// 按节点英文名过滤
		nid := logic.GetNidByEname(tab)
		if nid > 0 {
			topics = logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", "nid=? AND top!=1", nid)
			total = logic.DefaultTopic.Count(context.EchoContext(ctx), "nid=? AND top!=1", nid)
		} else {
			topics = logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", "top!=1")
			total = logic.DefaultTopic.Count(context.EchoContext(ctx), "top!=1")
		}
	}

	allTopics := append(topTopics, topics...)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"topics":   allTopics,
		"tab":      tab,
		"tab_list": hotNodes,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// WebsiteStat 网站统计信息
func (IndexController) WebsiteStat(ctx echo.Context) error {
	data := map[string]interface{}{
		"article":  logic.DefaultArticle.Total(),
		"project":  logic.DefaultProject.Total(),
		"topic":    logic.DefaultTopic.Total(),
		"resource": logic.DefaultResource.Total(),
		"book":     logic.DefaultGoBook.Total(),
		"comment":  logic.DefaultComment.Total(),
		"user":     logic.DefaultUser.Total(),
	}
	return success(ctx, data)
}
