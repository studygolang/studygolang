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

type TopicController struct{}

// RegisterRoute 注册路由
// 注意：node 路由必须在 :tid 参数路由前注册，避免路由冲突
func (self TopicController) RegisterRoute(g *echo.Group) {
	g.GET("/topics/no_reply", self.NoReply)
	g.GET("/topics/last", self.Last)
	g.GET("/topics/node/:nid", self.NodeTopics)
	g.GET("/topics", self.List)
	g.GET("/topics/:tid", self.Detail)
	g.GET("/nodes", self.Nodes)
}

// List 话题列表，支持 tab、p 参数
func (self TopicController) List(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab != "" && tab != "all" {
		nid := logic.GetNidByEname(tab)
		if nid > 0 {
			return self.topicList(ctx, tab, "topics.mtime DESC", "nid=? AND top!=1", nid)
		}
	}
	return self.topicList(ctx, "all", "topics.mtime DESC", "top!=1")
}

// NoReply 无回复话题列表
func (self TopicController) NoReply(ctx echo.Context) error {
	return self.topicList(ctx, "no_reply", "topics.mtime DESC", "lastreplyuid=?", 0)
}

// Last 最新话题列表
func (self TopicController) Last(ctx echo.Context) error {
	return self.topicList(ctx, "last", "ctime DESC", "")
}

func (TopicController) topicList(ctx echo.Context, tab, orderBy, querystring string, args ...interface{}) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	topTopics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "ctime DESC", "top=1")
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, orderBy, querystring, args...)
	total := logic.DefaultTopic.Count(context.EchoContext(ctx), querystring, args...)
	hasMore := paginator.SetTotal(total).HasMorePage()

	hotNodes := logic.DefaultTopic.FindHotNodes(context.EchoContext(ctx))

	return success(ctx, map[string]interface{}{
		"list":     append(topTopics, topics...),
		"tab":      tab,
		"tab_list": hotNodes,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// NodeTopics 节点下的话题列表
func (TopicController) NodeTopics(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	nid := goutils.MustInt(ctx.Param("nid"))
	querystring := "nid=?"
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", querystring, nid)
	total := logic.DefaultTopic.Count(context.EchoContext(ctx), querystring, nid)
	hasMore := paginator.SetTotal(total).HasMorePage()

	node := logic.GetNode(nid)

	return success(ctx, map[string]interface{}{
		"list":     topics,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
		"node":     node,
	})
}

// Detail 话题详情，增加浏览量
func (TopicController) Detail(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	topic, replies, err := logic.DefaultTopic.FindByTid(context.EchoContext(ctx), tid)
	if err != nil {
		return fail(ctx, "服务器异常")
	}

	me, ok := ctx.Get("user").(*model.Me)

	permission := topic["permission"].(int)
	switch permission {
	case model.PermissionLogin:
		if !ok {
			topic["content"] = "登录用户可见！"
		}
	case model.PermissionPay:
		if !ok || !me.IsVip || !me.IsRoot {
			topic["content"] = "付费用户可见！"
		}
	}

	logic.Views.Incr(Request(ctx), model.TypeTopic, tid)

	return success(ctx, map[string]interface{}{
		"topic":   topic,
		"replies": replies,
	})
}

// Nodes 获取所有节点列表
func (TopicController) Nodes(ctx echo.Context) error {
	nodes := logic.GenNodes()
	return success(ctx, nodes)
}
