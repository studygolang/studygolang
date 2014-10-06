// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"html/template"
	"model"
	"net/http"
	"service"
	"strconv"
	"util"
)

// 在需要评论且要回调的地方注册评论对象
func init() {
	// 注册评论对象
	service.RegisterCommentObject(model.TYPE_TOPIC, service.TopicComment{})
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
	switch vars["view"] {
	case "/no_reply":
		where = "lastreplyuid=0"
	case "/last":
		order = "ctime DESC"
	}
	topics, total := service.FindTopics(page, 0, where, order)
	pageHtml := service.GetPageHtml(page, total, "/topics")
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "nodes": nodes})
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
	}

	service.Views.Incr(req, model.TYPE_TOPIC, util.MustInt(vars["tid"]))

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/detail.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topic": topic, "replies": replies})
}

// 新建帖子
// uri: /topics/new{json:(|.json)}
func NewTopicHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.GenNodes()
	vars := mux.Vars(req)
	title := req.FormValue("title")
	// 请求新建帖子页面
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/new.html")
		filter.SetData(req, map[string]interface{}{"nodes": nodes})
		return
	}

	user, _ := filter.CurrentUser(req)
	// 入库
	topic := model.NewTopic()
	topic.Uid = user["uid"].(int)
	topic.Nid = util.MustInt(req.FormValue("nid"))
	topic.Title = req.FormValue("title")
	topic.Content = req.FormValue("content")
	errMsg, err := service.PublishTopic(topic)
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error":"`, errMsg, `"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}
