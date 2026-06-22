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

type IndexController struct{}

// RegisterRoute 注册路由
func (self IndexController) RegisterRoute(g *echo.Group) {
	g.GET("/home", self.Home)
	g.GET("/stat/site", self.WebsiteStat)
}

// Home 首页动态列表，使用 feed 表聚合多种类型内容
// 支持 tab 参数：all（默认，feed 动态）、recommend（推荐）、no_reply（未回复话题）
// 返回格式：feeds/tab/tab_list/total/page/has_more
func (IndexController) Home(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = model.TabAll
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	hotNodes := logic.DefaultTopic.FindHotNodes(context.EchoContext(ctx))

	// 默认使用 feed 聚合动态（话题、文章、项目、资源等）
	if tab == model.TabAll || tab == model.TabRecommend {
		// 获取置顶动态
		topFeeds := logic.DefaultFeed.FindTop(context.EchoContext(ctx))
		// 获取最新动态（带分页）
		feeds := logic.DefaultFeed.FindRecentWithPaginator(context.EchoContext(ctx), paginator, tab)
		// 获取总数
		total := logic.DefaultFeed.GetTotalCount(context.EchoContext(ctx))

		allFeeds := append(topFeeds, feeds...)
		hasMore := paginator.SetTotal(total).HasMorePage()

		// Email 脱敏：Feed.User / Lastreplyuser 自动序列化会暴露 Email
		sanitizeFeedsForPublic(allFeeds)

		return success(ctx, map[string]interface{}{
			"feeds":    allFeeds,
			"tab":      tab,
			"tab_list": hotNodes,
			"total":    total,
			"page":     curPage,
			"has_more": hasMore,
		})
	}

	// 其他 tab（如 no_reply、节点筛选）仍然返回话题列表
	var (
		topics []map[string]interface{}
		total  int64
	)

	switch tab {
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

	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"topics":   topics,
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
