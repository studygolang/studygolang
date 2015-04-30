// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"model"
	"service"
	"util"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	service.RegisterCommentObject(model.TYPE_TOPIC, service.TopicComment{})
	service.RegisterLikeObject(model.TYPE_TOPIC, service.TopicLike{})
}

// 社区帖子列表页
// uri: /topics{view:(|/popular|/no_reply|/last)}
func TopicsHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.GenNodes()
	// 设置内容模板
	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}
	vars := mux.Vars(req)
	order := ""
	where := ""
	view := ""
	switch vars["view"] {
	case "/no_reply":
		view = "no_reply"
		where = "lastreplyuid=0"
	case "/last":
		view = "last"
		order = "ctime DESC"
	}

	topics, total := service.FindTopics(page, 0, where, order)
	pageHtml := service.GetPageHtml(page, total, req.URL.Path)
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "nodes": nodes, "view": view})
}

// 某节点下的帖子列表
// uri: /topics/node{nid:[0-9]+}
func NodesHandler(rw http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}
	vars := mux.Vars(req)
	topics, total := service.FindTopics(page, 0, "nid="+vars["nid"])
	pageHtml := service.GetPageHtml(page, total, "/topics/node"+vars["nid"])
	// 当前节点信息
	node := service.GetNode(util.MustInt(vars["nid"]))
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/node.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "total": total, "node": node})
}

// 社区帖子详细页
// uri: /topics/{tid:[0-9]+}
func TopicDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	topic, replies, err := service.FindTopicByTid(vars["tid"])
	if err != nil {
		util.Redirect(rw, req, "/topics")
		return
	}

	likeFlag := 0
	hadCollect := 0
	user, ok := filter.CurrentUser(req)
	if ok {
		uid := user["uid"].(int)
		tid := topic["tid"].(int)
		likeFlag = service.HadLike(uid, tid, model.TYPE_TOPIC)
		hadCollect = service.HadFavorite(uid, tid, model.TYPE_TOPIC)
	}

	service.Views.Incr(req, model.TYPE_TOPIC, util.MustInt(vars["tid"]))

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/detail.html,/template/common/comment.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topic": topic, "replies": replies, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// 新建帖子
// uri: /topics/new{json:(|.json)}
func NewTopicHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.GenNodes()
	vars := mux.Vars(req)
	title := req.PostFormValue("title")
	// 请求新建主题页面
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/new.html")
		filter.SetData(req, map[string]interface{}{"nodes": nodes, "activeTopics": "active"})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishTopic(user, req.PostForm)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}

// 修改主题
// uri: /topics/modify{json:(|.json)}
func ModifyTopicHandler(rw http.ResponseWriter, req *http.Request) {
	tid := req.FormValue("tid")
	if tid == "" {
		util.Redirect(rw, req, "/topics")
		return
	}

	nodes := service.GenNodes()

	vars := mux.Vars(req)
	// 请求编辑主题页面
	if req.Method != "POST" || vars["json"] == "" {
		topic := service.FindTopic(tid)
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/new.html")
		filter.SetData(req, map[string]interface{}{"nodes": nodes, "topic": topic, "activeTopics": "active"})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishTopic(user, req.PostForm)
	if err != nil {
		if err == service.NotModifyAuthorityErr {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}
