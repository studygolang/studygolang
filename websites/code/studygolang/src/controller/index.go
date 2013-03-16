// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"net/http"
	"service"
)

// 首页
func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := genNodes()
	// 获取最新帖子
	newTopics, _ := service.FindTopics(1, 10, "", "ctime DESC")
	// 获取热门帖子
	hotTopics := service.FindHotTopics()
	// 获得最新博文
	articles := service.FindNewBlogs()
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/index.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"news": newTopics, "hots": hotTopics, "articles": articles, "nodes": nodes})
}
