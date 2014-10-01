// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"model"
	"net/http"
	"service"
	"strconv"
)

// 侧边栏的内容通过异步请求获取

// 某节点下其他帖子
func OtherTopicsHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	topics := service.FindTopicsByNid(vars["nid"], vars["tid"])
	topics = service.JSEscape(topics)
	data, err := json.Marshal(topics)
	if err != nil {
		logger.Errorln("[OtherTopicsHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "topics":`+string(data)+`}`)
}

// 网站统计信息
func StatHandler(rw http.ResponseWriter, req *http.Request) {
	topicTotal := service.TopicsTotal()
	replyTotal := service.CommentsTotal(model.TYPE_TOPIC)
	resourceTotal := service.ResourcesTotal()
	userTotal := service.CountUsers()
	fmt.Fprint(rw, `{"errno": 0, "topic":`+strconv.Itoa(topicTotal)+`,"resource":`+strconv.Itoa(resourceTotal)+`,"reply":`+strconv.Itoa(replyTotal)+`,"user":`+strconv.Itoa(userTotal)+`}`)
}

// 社区最新公告
// uri: /topics/notice.json
func NoticeHandler(rw http.ResponseWriter, req *http.Request) {
	topic := service.FindNoticeTopic()
	newNotice, err := json.Marshal(topic)
	if err != nil {
		logger.Errorln("[NoticeHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "notice":`+string(newNotice)+`}`)
}

// 最新帖子
// uri: /topics/recent.json
func RecentTopicHandler(rw http.ResponseWriter, req *http.Request) {
	recentTopics := service.FindRecentTopics(0, "10")
	buf, err := json.Marshal(recentTopics)
	if err != nil {
		logger.Errorln("[RecentTopicHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 最新博文
// uri: /articles/recent.json
func RecentArticleHandler(rw http.ResponseWriter, req *http.Request) {
	recentArticles := service.FindArticles("0", "10")
	buf, err := json.Marshal(recentArticles)
	if err != nil {
		logger.Errorln("[RecentArticleHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 最新评论
// uri: /comments/recent.json
func RecentCommentHandler(rw http.ResponseWriter, req *http.Request) {
	recentComments := service.FindRecentComments(0, -1, "10")
	buf, err := json.Marshal(recentComments)
	if err != nil {
		logger.Errorln("[RecentArticleHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 社区热门节点
// uri: /nodes/hot.json
func HotNodesHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.FindHotNodes()
	hotNodes, err := json.Marshal(nodes)
	if err != nil {
		logger.Errorln("[HotNodesHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "nodes":`+string(hotNodes)+`}`)
}
