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

// Home 首页动态列表，优先按后台 IndexNavs 配置返回对应数据源。
func (IndexController) Home(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab == "" && len(logic.WebsiteSetting.IndexNavs) > 0 {
		tab = logic.WebsiteSetting.IndexNavs[0].Tab
	}
	if tab == "" {
		tab = model.TabAll
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	hotNodes := logic.DefaultTopic.FindHotNodes(context.EchoContext(ctx))

	if len(logic.WebsiteSetting.IndexNavs) > 0 {
		data := logic.DefaultIndex.FindData(context.EchoContext(ctx), tab, paginator)
		sanitizeHomeDataForPublic(data)
		data["tab_list"] = hotNodes
		data["page"] = curPage
		data["has_more"] = false

		currentTab := data["tab"]
		if currentTab == model.TabAll || currentTab == model.TabRecommend {
			total := logic.DefaultFeed.GetTotalCount(context.EchoContext(ctx))
			data["total"] = total
			data["has_more"] = paginator.SetTotal(total).HasMorePage()
		}

		return success(ctx, data)
	}

	// 未配置 IndexNavs 时保留原首页 API 行为。
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

func sanitizeHomeDataForPublic(data map[string]interface{}) {
	if feeds, ok := data["feeds"].([]*model.Feed); ok {
		sanitizeFeedsForPublic(feeds)
	}
	if topics, ok := data["topics"].([]map[string]interface{}); ok {
		for _, topic := range topics {
			if user, ok := topic["user"].(*model.User); ok {
				sanitizeUserForPublic(user)
			}
		}
	}
	if articles, ok := data["articles"].([]*model.Article); ok {
		for _, article := range articles {
			sanitizeUserForPublic(article.LastReplyUser)
		}
	}
	if users, ok := data["users"].(map[int]*model.User); ok {
		for _, user := range users {
			sanitizeUserForPublic(user)
		}
	}
	if subjects, ok := data["subjects"].([]*model.Subject); ok {
		sanitizeSubjectsForPublic(subjects)
	}
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
