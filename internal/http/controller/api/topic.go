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
	g.POST("/topics", self.Publish)
	g.GET("/topics/:tid/edit", self.Edit)
	g.PUT("/topics/:tid", self.Update)
	g.GET("/topics/:tid", self.Detail)
	g.POST("/topics/:tid/set_top", self.SetTop)
	g.POST("/topics/:tid/append", self.Append)
	g.GET("/topics/:tid/appends", self.Appends)
	g.GET("/topics/:nid/others", self.OthersTopics)
	g.GET("/nodes", self.Nodes)
}

// List 话题列表，支持 tab、p、sort 参数
func (self TopicController) List(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	sort := ctx.QueryParam("sort")

	orderBy := getTopicSortOrder(sort)

	if tab != "" && tab != "all" {
		nid := logic.GetNidByEname(tab)
		if nid > 0 {
			return self.topicList(ctx, tab, orderBy, "nid=? AND top!=1", nid)
		}
	}
	return self.topicList(ctx, "all", orderBy, "top!=1")
}

// getTopicSortOrder 根据 sort 参数返回排序 SQL
// hot: 按回复数+点赞数排序（综合热度）
// latest: 按 id DESC（最新发布）
// noreply: 按回复数 ASC（无人回复优先）
func getTopicSortOrder(sort string) string {
	switch sort {
	case "hot":
		// 回复数*2 + 点赞数，发布时间越近权重越高
		return "topics.reply+topics.`like`*2 DESC, topics.id DESC"
	case "latest":
		return "topics.id DESC"
	case "noreply":
		return "topics.reply ASC, topics.id DESC"
	default:
		// 默认按最近活动时间
		return "topics.mtime DESC"
	}
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
	if curPage < 1 {
		curPage = 1
	}
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

	// 防御性检查：确保 topic 不为 nil
	if topic == nil {
		return fail(ctx, "话题不存在")
	}

	me, ok := ctx.Get("user").(*model.Me)

	permission, _ := topic["permission"].(int)
	uid, _ := topic["uid"].(int)

	switch permission {
	case model.PermissionOnlyMe:
		// 仅作者和管理员可见
		if !ok || (uid != me.Uid && !me.IsRoot) {
			return fail(ctx, "话题不存在")
		}
	case model.PermissionLogin:
		if !ok {
			topic["content"] = "登录用户可见！"
		}
	case model.PermissionPay:
		if !ok || (!me.IsVip && !me.IsRoot && uid != me.Uid) {
			topic["content"] = "付费用户可见！"
		}
	}

	// 已登录用户的附加信息
	result := map[string]interface{}{
		"topic":   topic,
		"replies": replies,
	}

	if ok {
		tidInt, _ := topic["tid"].(int)
		result["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, tidInt, model.TypeTopic)
		result["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, tidInt, model.TypeTopic)

		logic.Views.Incr(Request(ctx), model.TypeTopic, tidInt, me.Uid)
	} else {
		logic.Views.Incr(Request(ctx), model.TypeTopic, tid)
	}

	// 附言：公开、登录可见且已登录、付费且已满足条件
	canViewAppends := permission == model.PermissionPublic ||
		(permission == model.PermissionLogin && ok) ||
		(permission == model.PermissionPay && ok && (me.IsVip || me.IsRoot || uid == me.Uid)) ||
		permission == model.PermissionOnlyMe
	if canViewAppends {
		result["appends"] = logic.DefaultTopic.FindAppend(context.EchoContext(ctx), tid)
	}

	return success(ctx, result)
}

// Edit 获取话题编辑数据（需要登录，验证权限）
func (TopicController) Edit(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	topics := logic.DefaultTopic.FindByTids([]int{tid})
	if len(topics) == 0 {
		return fail(ctx, "话题不存在")
	}

	topic := topics[0]
	// 验证权限：只能编辑自己的话题，或管理员可以编辑所有话题
	if !logic.CanEdit(me, topic) {
		return fail(ctx, "没有编辑权限")
	}

	// 获取节点列表
	nodes := logic.GenNodes()

	return success(ctx, map[string]interface{}{
		"topic": topic,
		"nodes": nodes,
	})
}

// Update 更新话题（需要登录，验证权限）
func (TopicController) Update(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	topics := logic.DefaultTopic.FindByTids([]int{tid})
	if len(topics) == 0 {
		return fail(ctx, "话题不存在")
	}

	topic := topics[0]
	// 验证权限
	if !logic.CanEdit(me, topic) {
		return fail(ctx, "没有编辑权限")
	}

	forms, _ := ctx.FormParams()
	forms.Set("tid", ctx.Param("tid"))
	errMsg, err := logic.DefaultTopic.Modify(context.EchoContext(ctx), me, forms)
	if err != nil {
		if errMsg != "" {
			return fail(ctx, errMsg)
		}
		return fail(ctx, "更新失败")
	}

	return success(ctx, map[string]interface{}{"tid": tid})
}

// OthersTopics 获取同一节点下的其他话题（用于话题详情页侧边栏）
// GET /api/v1/topics/:nid/others?limit=10
func (TopicController) OthersTopics(ctx echo.Context) error {
	nid := goutils.MustInt(ctx.Param("nid"), 0)
	if nid == 0 {
		return fail(ctx, "节点 ID 不合法")
	}

	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	curPage := 1
	paginator := logic.NewPaginatorWithPerPage(curPage, limit)
	querystring := "nid=?"
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", querystring, nid)
	node := logic.GetNode(nid)

	return success(ctx, map[string]interface{}{
		"list": topics,
		"node": node,
	})
}

// Nodes 获取所有节点列表
// GET /api/v1/nodes
func (TopicController) Nodes(ctx echo.Context) error {
	data := make(map[string]interface{})
	if len(logic.AllRecommendNodes) > 0 {
		data["nodes"] = logic.DefaultNode.FindAll(context.EchoContext(ctx))
	} else {
		data["nodes"] = logic.GenNodes()
	}
	return success(ctx, data)
}

// Publish 发布新话题（需要登录，统一使用 requireAuth 认证）
func (TopicController) Publish(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	forms, _ := ctx.FormParams()
	tid, err := logic.DefaultTopic.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败："+err.Error())
	}
	return success(ctx, map[string]interface{}{"tid": tid})
}

// Append 发布话题附言（需要登录 + 作者权限校验）
func (TopicController) Append(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	// 验证话题是否存在及权限
	topics := logic.DefaultTopic.FindByTids([]int{tid})
	if len(topics) == 0 {
		return fail(ctx, "话题不存在")
	}
	if topics[0].Uid != me.Uid && !me.IsRoot {
		return fail(ctx, "只有话题作者才能添加附言")
	}

	content := ctx.FormValue("content")
	if content == "" {
		return fail(ctx, "附言内容不能为空")
	}
	if len(content) > 65535 {
		return fail(ctx, "附言内容过长")
	}

	err = logic.DefaultTopic.Append(context.EchoContext(ctx), me.Uid, tid, content)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{"tid": tid})
}

// Appends 获取话题附言列表
func (TopicController) Appends(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	appends := logic.DefaultTopic.FindAppend(context.EchoContext(ctx), tid)
	return success(ctx, map[string]interface{}{
		"appends": appends,
	})
}

// SetTop 设置话题置顶（需要登录）
func (TopicController) SetTop(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return fail(ctx, "tid 非法")
	}

	err = logic.DefaultTopic.SetTop(context.EchoContext(ctx), me, tid)
	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{"tid": tid})
}
