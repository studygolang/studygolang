// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"http/middleware"
	"logic"
	"model"
	"net/http"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeTopic, logic.TopicComment{})
	logic.RegisterLikeObject(model.TypeTopic, logic.TopicLike{})
}

type TopicController struct{}

// 注册路由
func (self TopicController) RegisterRoute(g *echo.Group) {
	g.GET("/topics", self.Topics)
	g.GET("/topics/no_reply", self.TopicsNoReply)
	g.GET("/topics/last", self.TopicsLast)
	g.GET("/topics/:tid", self.Detail)
	g.GET("/topics/node/:nid", self.NodeTopics)

	g.Match([]string{"GET", "POST"}, "/topics/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.PublishNotice())
	g.Match([]string{"GET", "POST"}, "/topics/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
}

func (self TopicController) Topics(ctx echo.Context) error {
	return self.topicList(ctx, "", "topics.mtime DESC", "")
}

func (self TopicController) TopicsNoReply(ctx echo.Context) error {
	return self.topicList(ctx, "no_reply", "topics.mtime DESC", "lastreplyuid=?", 0)
}

func (self TopicController) TopicsLast(ctx echo.Context) error {
	return self.topicList(ctx, "last", "ctime DESC", "")
}

func (TopicController) topicList(ctx echo.Context, view, orderBy, querystring string, args ...interface{}) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	topics := logic.DefaultTopic.FindAll(ctx, paginator, orderBy, querystring, args...)
	total := logic.DefaultTopic.Count(ctx, querystring, args...)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	data := map[string]interface{}{
		"topics":       topics,
		"activeTopics": "active",
		"nodes":        logic.GenNodes(),
		"view":         view,
		"page":         template.HTML(pageHtml),
	}

	return render(ctx, "topics/list.html", data)
}

// NodeTopics 某节点下的主题列表
func (TopicController) NodeTopics(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring, nid := "nid=?", goutils.MustInt(ctx.Param("nid"))
	topics := logic.DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nid)
	total := logic.DefaultTopic.Count(ctx, querystring, nid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	// 当前节点信息
	node := logic.GetNode(nid)

	return render(ctx, "topics/node.html", map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "total": total, "node": node})
}

// Detail 社区主题详细页
func (TopicController) Detail(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	topic, replies, err := logic.DefaultTopic.FindByTid(ctx, tid)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	likeFlag := 0
	hadCollect := 0
	user, ok := ctx.Get("user").(*model.Me)
	if ok {
		tid := topic["tid"].(int)
		likeFlag = logic.DefaultLike.HadLike(ctx, user.Uid, tid, model.TypeTopic)
		hadCollect = logic.DefaultFavorite.HadFavorite(ctx, user.Uid, tid, model.TypeTopic)
	}

	logic.Views.Incr(Request(ctx), model.TypeTopic, tid)

	return render(ctx, "topics/detail.html,common/comment.html", map[string]interface{}{"activeTopics": "active", "topic": topic, "replies": replies, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// Create 新建主题
func (TopicController) Create(ctx echo.Context) error {
	nodes := logic.GenNodes()

	title := ctx.FormValue("title")
	// 请求新建主题页面
	if title == "" || ctx.Request().Method() != "POST" {
		return render(ctx, "topics/new.html", map[string]interface{}{"nodes": nodes, "activeTopics": "active"})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultTopic.Publish(ctx, me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 1, "内部服务错误")
	}

	return success(ctx, nil)
}

// Modify 修改主题
func (TopicController) Modify(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.FormValue("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	nodes := logic.GenNodes()

	if ctx.Request().Method() != "POST" {
		topics := logic.DefaultTopic.FindByTids([]int{tid})
		if len(topics) == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/topics")
		}

		return render(ctx, "topics/new.html", map[string]interface{}{"nodes": nodes, "topic": topics[0], "activeTopics": "active"})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultTopic.Publish(ctx, me, ctx.FormParams())
	if err != nil {
		if err == logic.NotModifyAuthorityErr {
			return fail(ctx, 1, "没有权限操作")
		}

		return fail(ctx, 2, "服务错误，请稍后重试！")
	}
	return success(ctx, nil)
}
