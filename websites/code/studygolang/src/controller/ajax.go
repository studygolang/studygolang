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
