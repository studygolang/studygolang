// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"net/http"
	"net/url"
	"strings"

	"config"
	"filter"
	"logger"
	"model"
	"service"
	"util"
)

// 首页
func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	// nodes := service.GenNodes()

	num := 10

	topicsList := make([]map[string]interface{}, num)
	// 置顶的topic
	topTopics, _ := service.FindTopics(1, num, "top=1", "ctime DESC")
	if len(topTopics) < num {
		// 获取最新帖子
		newTopics, _ := service.FindTopics(1, num-len(topTopics), "top=0", "ctime DESC")

		topicsList = append(topTopics, newTopics...)
	}

	// 获取热门帖子
	//hotTopics := service.FindHotTopics()
	// 获得最新博文
	// blogs := service.FindNewBlogs()
	recentArticles := service.FindArticles("0", "10")
	// 获取当前用户喜欢对象信息
	var likeFlags map[int]int

	if len(recentArticles) > 0 {
		user, ok := filter.CurrentUser(req)
		if ok {
			uid := user["uid"].(int)

			likeFlags, _ = service.FindUserLikeObjects(uid, model.TYPE_ARTICLE, recentArticles[0].Id, recentArticles[len(recentArticles)-1].Id)
		}
	}

	// Golang 资源
	resources := service.FindResources("0", "10")

	/*
		start, end := 0, len(resources)
		if n := end - 10; n > 0 {
			start = rand.Intn(n)
			end = start + 10
		}
	*/

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/index.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"topics": topicsList, "articles": recentArticles, "likeflags": likeFlags, "resources": resources})
}

// 包装链接
func WRHandler(rw http.ResponseWriter, req *http.Request) {
	tUrl := req.FormValue("u")
	if tUrl == "" {
		util.Redirect(rw, req, "/")
		return
	}

	if pUrl, err := url.Parse(tUrl); err != nil {
		util.Redirect(rw, req, tUrl)
		return
	} else {
		if !pUrl.IsAbs() {
			util.Redirect(rw, req, tUrl)
			return
		}

		// 本站
		if strings.Contains(pUrl.Host, config.Config["domain"]) {
			util.Redirect(rw, req, tUrl)
			return
		}

		// 检测是否禁止了 iframe 加载
		// 看是否在黑名单中
		for _, denyHost := range strings.Split(config.Config["iframe_deny"], ",") {
			if strings.Contains(pUrl.Host, denyHost) {
				util.Redirect(rw, req, tUrl)
				return
			}
		}

		// 检测会比较慢，进行异步检测，记录下来，以后分析再加黑名单
		go func() {
			resp, err := http.Head(tUrl)
			if err != nil {
				logger.Errorln("[iframe] head url:", tUrl, "error:", err)
				return
			}
			defer resp.Body.Close()
			if resp.Header.Get("X-Frame-Options") != "" {
				logger.Errorln("[iframe] deny:", tUrl)
				return
			}
		}()
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/wr.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"url": tUrl})
}

// PkgdocHandler Go 语言文档中文版
func PkgdocHandler(rw http.ResponseWriter, req *http.Request) {
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/pkgdoc.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeDoc": "active"})
}
