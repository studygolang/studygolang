// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/studygolang/mux"
	"logger"
	"service"
	"util"
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
// uri: /websites/stat.json
func StatHandler(rw http.ResponseWriter, req *http.Request) {
	articleTotal := service.ArticlesTotal()
	topicTotal := service.TopicsTotal()
	cmtTotal := service.CommentsTotal(-1)
	resourceTotal := service.ResourcesTotal()
	userTotal := service.CountUsers()

	data := map[string]int{
		"article":  articleTotal,
		"topic":    topicTotal,
		"resource": resourceTotal,
		"comment":  cmtTotal,
		"user":     userTotal,
	}

	buf, err := json.Marshal(data)
	if err != nil {
		logger.Errorln("[StatHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 社区最新公告或go最新动态
// uri: /dymanics/recent.json
func RecentDymanicHandler(rw http.ResponseWriter, req *http.Request) {
	dynamics := service.FindDynamics("0", "3")
	buf, err := json.Marshal(dynamics)
	if err != nil {
		logger.Errorln("[RecentDymanicHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
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

// 最新资源
// uri: /resources/recent.json
func RecentResourceHandler(rw http.ResponseWriter, req *http.Request) {
	recentResources := service.FindResources("0", "10")
	buf, err := json.Marshal(recentResources)
	if err != nil {
		logger.Errorln("[RecentResourceHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 最新评论
// uri: /comments/recent.json
func RecentCommentHandler(rw http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "10"
	}
	recentComments := service.FindRecentComments(0, -1, limit)

	uids := util.Models2Intslice(recentComments, "Uid")
	users := service.GetUserInfos(uids)

	result := map[string]interface{}{
		"comments": recentComments,
	}

	// json encode 不支持 map[int]...
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	buf, err := json.Marshal(result)

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

// 活跃会员
// uri: /user/active.json
func ActiveUserHandler(rw http.ResponseWriter, req *http.Request) {
	activeUsers := service.FindActiveUsers(0, 9)
	buf, err := json.Marshal(activeUsers)
	if err != nil {
		logger.Errorln("[ActiveUserHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

const maxImageSize = 5 << 20 // 5M

func UploadImageHandler(rw http.ResponseWriter, req *http.Request) {
	file, fileHeader, err := req.FormFile("img")
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"非法文件上传！"}`)
		return
	}

	// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
	if _, ok := file.(*os.File); ok {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
		return
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件读取失败！"}`)
		return
	}

	if len(buf) > maxImageSize {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
		return
	}

	uri := util.DateNow() + "/" + util.Md5Buf(buf) + filepath.Ext(fileHeader.Filename)

	err = service.UploadMemoryFile(file, uri)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件上传失败！"}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "uri":"`+uri+`"}`)
}
