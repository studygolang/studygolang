// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"logger"
	"model"
	"service"
	"util"

	"github.com/studygolang/mux"
)

// 侧边栏的内容通过异步请求获取

// 技术晨读
// uri: /readings/recent.json
func RecentReadingHandler(rw http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "7"
	}

	readings := service.FindReadings("0", limit, model.RtypeGo)
	buf, err := json.Marshal(readings)
	if err != nil {
		logger.Errorln("[RecentReadingHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

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
	projectTotal := service.ProjectsTotal()
	topicTotal := service.TopicsTotal()
	cmtTotal := service.CommentsTotal(-1)
	resourceTotal := service.ResourcesTotal()
	userTotal := service.CountUsers()

	data := map[string]int{
		"article":  articleTotal,
		"project":  projectTotal,
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
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "10"
	}

	recentTopics := service.FindRecentTopics(0, limit)
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
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "10"
	}

	recentArticles := service.FindArticles("0", limit)
	buf, err := json.Marshal(recentArticles)
	if err != nil {
		logger.Errorln("[RecentArticleHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 最新开源项目（TODO：暂时打乱，避免每次一样）
// uri: /projects/recent.json
func RecentProjectHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		limit = 10
		err   error
	)
	limitStr := req.FormValue("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			fmt.Fprint(rw, `{"ok": 0, "error":"limit is not number"}`)
			return
		}
	}

	recentProjects := service.FindProjects("0", "100")
	start, end := 0, len(recentProjects)
	if n := end - limit; n > 0 {
		start = rand.Intn(n)
		end = start + limit
	}

	buf, err := json.Marshal(recentProjects[start:end])
	if err != nil {
		logger.Errorln("[RecentProjectHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 最新资源
// uri: /resources/recent.json
func RecentResourceHandler(rw http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "10"
	}

	recentResources := service.FindResources("0", limit)
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
		logger.Errorln("[RecentCommentHandler] json.marshal error:", err)
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
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(hotNodes)+`}`)
}

// 活跃会员
// uri: /users/active.json
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

// 新加入会员
// uri: /users/newest.json
func NewestUserHandler(rw http.ResponseWriter, req *http.Request) {
	newestUsers := service.FindNewUsers(0, 9)
	buf, err := json.Marshal(newestUsers)
	if err != nil {
		logger.Errorln("[NewestUserHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 评论或回复 @ 某人 suggest
// uri: /at/users.json
func AtUsersHandler(rw http.ResponseWriter, req *http.Request) {
	term := req.FormValue("term")
	users := service.GetUserMentions(term, 10)
	buf, err := json.Marshal(users)
	if err != nil {
		logger.Errorln("[AtUsersHandler] json.marshal error:", err)
		fmt.Fprint(rw, `[]`)
		return
	}

	fmt.Fprint(rw, string(buf))
}

const maxImageSize = 5 << 20 // 5M

func UploadImageHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		uri    string
		buf    []byte
		err    error
		reader io.Reader
	)

	origUrl := req.FormValue("url")
	if origUrl != "" {
		resp, err := http.Get(origUrl)
		if err != nil {
			fmt.Fprint(rw, `{"ok": 0, "error":"获取图片失败"}`)
			return
		}
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)

		ext := filepath.Ext(origUrl)
		if ext == "" {
			contentType := http.DetectContentType(buf)
			exts, _ := mime.ExtensionsByType(contentType)
			if len(exts) > 0 {
				ext = exts[0]
			}
		}

		uri = util.DateNow() + "/" + util.Md5Buf(buf) + ext

		reader = bytes.NewReader(buf)
	} else {

		file, fileHeader, err := req.FormFile("img")
		if err != nil {
			fmt.Fprint(rw, `{"ok": 0, "error":"非法文件上传！"}`)
			return
		}

		defer file.Close()

		// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
		if _, ok := file.(*os.File); ok {
			fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
			return
		}

		reader = file

		buf, err := ioutil.ReadAll(file)
		imgDir := util.DateNow()
		if req.FormValue("avatar") != "" {
			imgDir = "avatar"
		}
		uri = imgDir + "/" + util.Md5Buf(buf) + filepath.Ext(fileHeader.Filename)
	}

	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件读取失败！"}`)
		return
	}

	if len(buf) > maxImageSize {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
		return
	}

	err = service.UploadMemoryFile(reader, uri)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件上传失败！"}`)
		return
	}

	if origUrl != "" {
		uri = "http://studygolang.qiniudn.com/" + uri
	}

	fmt.Fprint(rw, `{"ok": 1, "uri":"`+uri+`"}`)
}
