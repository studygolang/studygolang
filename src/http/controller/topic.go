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
	"strconv"
	"util"

	. "http"

	"github.com/dchest/captcha"
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
	g.GET("/topics", self.TopicList)
	g.GET("/topics/no_reply", self.TopicsNoReply)
	g.GET("/topics/last", self.TopicsLast)
	g.GET("/topics/:tid", self.Detail)
	g.GET("/topics/node/:nid", self.NodeTopics)
	g.GET("/go/:node", self.GoNodeTopics)
	g.GET("/nodes", self.Nodes)

	g.Match([]string{"GET", "POST"}, "/topics/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice(), middleware.CheckCaptcha())
	g.Match([]string{"GET", "POST"}, "/topics/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())

	g.POST("/topics/set_top", self.SetTop, middleware.NeedLogin())

	g.Match([]string{"GET", "POST"}, "/append/topic/:tid", self.Append, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck())
}

func (self TopicController) TopicList(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = GetFromCookie(ctx, "TOPIC_TAB")
	} else {
		SetCookie(ctx, "TOPIC_TAB", tab)
	}

	if tab != "" && tab != "all" {
		nid := logic.GetNidByEname(tab)
		if nid > 0 {
			return self.topicList(ctx, tab, "topics.mtime DESC", "nid=? AND top!=1", nid)
		}
	}

	return self.topicList(ctx, "all", "topics.mtime DESC", "top!=1")
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

func (TopicController) topicList(ctx echo.Context, tab, orderBy, querystring string, args ...interface{}) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	// 置顶的topic
	topTopics := logic.DefaultTopic.FindAll(ctx, paginator, "ctime DESC", "top=1")

	topics := logic.DefaultTopic.FindAll(ctx, paginator, orderBy, querystring, args...)
	total := logic.DefaultTopic.Count(ctx, querystring, args...)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	hotNodes := logic.DefaultTopic.FindHotNodes(ctx)

	data := map[string]interface{}{
		"topics":       append(topTopics, topics...),
		"activeTopics": "active",
		"nodes":        logic.GenNodes(),
		"tab":          tab,
		"tab_list":     hotNodes,
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

// GoNodeTopics 某节点下的主题列表，uri: /go/golang
func (TopicController) GoNodeTopics(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	ename := ctx.Param("node")
	node := logic.GetNodeByEname(ename)
	if node == nil {
		return render(ctx, "notfound.html", nil)
	}

	querystring, nid := "nid=?", node["nid"].(int)
	topics := logic.DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nid)
	total := logic.DefaultTopic.Count(ctx, querystring, nid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	return render(ctx, "topics/node.html", map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "total": total, "node": node})
}

// Detail 社区主题详细页
func (TopicController) Detail(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return render(ctx, "notfound.html", nil)
	}

	topic, replies, err := logic.DefaultTopic.FindByTid(ctx, tid)
	if err != nil {
		return render(ctx, "notfound.html", nil)
	}

	data := map[string]interface{}{
		"activeTopics": "active",
		"topic":        topic,
		"replies":      replies,
		"appends":      []*model.TopicAppend{},
	}

	me, ok := ctx.Get("user").(*model.Me)
	if topic["permission"] == 0 || (topic["permission"] == 1 && ok) {
		data["appends"] = logic.DefaultTopic.FindAppend(ctx, tid)
	}
	if ok {
		tid := topic["tid"].(int)
		data["likeflag"] = logic.DefaultLike.HadLike(ctx, me.Uid, tid, model.TypeTopic)
		data["hadcollect"] = logic.DefaultFavorite.HadFavorite(ctx, me.Uid, tid, model.TypeTopic)

		logic.Views.Incr(Request(ctx), model.TypeTopic, tid, me.Uid)

		if me.Uid != topic["uid"].(int) {
			go logic.DefaultViewRecord.Record(tid, model.TypeTopic, me.Uid)
		}

		if me.IsRoot || me.Uid == topic["uid"].(int) {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(ctx, tid, model.TypeTopic)
			data["view_source"] = logic.DefaultViewSource.FindOne(ctx, tid, model.TypeTopic)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeTopic, tid)
	}

	return render(ctx, "topics/detail.html,common/comment.html", data)
}

// Create 新建主题
func (TopicController) Create(ctx echo.Context) error {
	nid := goutils.MustInt(ctx.FormValue("nid"))

	me := ctx.Get("user").(*model.Me)

	title := ctx.FormValue("title")
	// 请求新建主题页面
	if title == "" || ctx.Request().Method() != "POST" {
		hotNodes := logic.DefaultTopic.FindHotNodes(ctx)

		data := map[string]interface{}{
			"activeTopics": "active",
			"nid":          nid,
			"tab_list":     hotNodes,
		}

		if logic.NeedCaptcha(me) {
			data["captchaId"] = captcha.NewLen(util.CaptchaLen)
		}

		hadRecommend := false
		if len(logic.AllRecommendNodes) > 0 {
			hadRecommend = true

			data["nodes"] = logic.DefaultNode.FindAll(ctx)
		} else {
			data["nodes"] = logic.GenNodes()
		}

		data["had_recommend"] = hadRecommend

		return render(ctx, "topics/new.html", data)
	}

	if nid == 0 {
		return fail(ctx, 1, "没有选择节点！")
	}

	tid, err := logic.DefaultTopic.Publish(ctx, me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 3, "内部服务错误:"+err.Error())
	}

	return success(ctx, map[string]interface{}{"tid": tid})
}

// Modify 修改主题
func (TopicController) Modify(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.FormValue("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	if ctx.Request().Method() != "POST" {
		topics := logic.DefaultTopic.FindByTids([]int{tid})
		if len(topics) == 0 {
			return ctx.Redirect(http.StatusSeeOther, "/topics")
		}

		hotNodes := logic.DefaultTopic.FindHotNodes(ctx)

		data := map[string]interface{}{
			"topic":        topics[0],
			"activeTopics": "active",
			"tab_list":     hotNodes,
		}

		hadRecommend := false
		if len(logic.AllRecommendNodes) > 0 {
			hadRecommend = true

			data["nodes"] = logic.DefaultNode.FindAll(ctx)
		} else {
			data["nodes"] = logic.GenNodes()
		}

		data["had_recommend"] = hadRecommend

		return render(ctx, "topics/new.html", data)
	}

	me := ctx.Get("user").(*model.Me)
	_, err := logic.DefaultTopic.Publish(ctx, me, ctx.FormParams())
	if err != nil {
		if err == logic.NotModifyAuthorityErr {
			return fail(ctx, 1, "没有权限操作")
		}

		return fail(ctx, 2, "服务错误，请稍后重试！")
	}
	return success(ctx, map[string]interface{}{"tid": tid})
}

func (TopicController) Append(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	topics := logic.DefaultTopic.FindByTids([]int{tid})
	if len(topics) == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	topic := topics[0]
	me := ctx.Get("user").(*model.Me)
	if topic.Uid != me.Uid {
		return ctx.Redirect(http.StatusSeeOther, "/topics/"+strconv.Itoa(tid))
	}

	// 请求新建主题页面
	if ctx.Request().Method() != http.MethodPost {
		data := map[string]interface{}{
			"topic":        topic,
			"activeTopics": "active",
		}

		return render(ctx, "topics/append.html", data)
	}

	content := ctx.FormValue("content")
	err := logic.DefaultTopic.Append(ctx, me.Uid, tid, content)
	if err != nil {
		return fail(ctx, 1, "出错了:"+err.Error())
	}

	return success(ctx, nil)
}

// Nodes 所有节点
func (TopicController) Nodes(ctx echo.Context) error {
	data := make(map[string]interface{})

	if len(logic.AllRecommendNodes) > 0 {
		data["nodes"] = logic.DefaultNode.FindAll(ctx)
	} else {
		data["nodes"] = logic.GenNodes()
	}

	return render(ctx, "topics/nodes.html", data)
}

func (TopicController) SetTop(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.FormValue("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultTopic.SetTop(ctx, me, tid)
	if err != nil {
		if err == logic.NotFoundErr {
			return ctx.Redirect(http.StatusSeeOther, "/topics")
		}

		return fail(ctx, 1, "出错了:"+err.Error())
	}

	return success(ctx, nil)
}
