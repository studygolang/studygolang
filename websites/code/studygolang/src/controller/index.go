// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"math/rand"
	"net/http"
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
	// TODO：开源项目（暂时使用 resource 表）
	resources := service.FindResourcesByCatid("2")

	start, end := 0, len(resources)
	if n := end - 10; n > 0 {
		start = rand.Intn(n)
		end = start + 10
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/index.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"topics": newTopics, "articles": recentArticles, "resources": resources[start:end]})
}
