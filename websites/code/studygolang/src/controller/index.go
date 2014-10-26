// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"math/rand"
	"net/http"

	"filter"
	"model"
	"service"
)

// 首页
func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	// nodes := service.GenNodes()

	// 获取最新帖子
	newTopics, _ := service.FindTopics(1, 10, "", "ctime DESC")
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
	resources := service.FindResources("0", "100")

	start, end := 0, len(resources)
	if n := end - 10; n > 0 {
		start = rand.Intn(n)
		end = start + 10
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/index.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"topics": newTopics, "articles": recentArticles, "likeflags": likeFlags, "resources": resources[start:end]})
}
