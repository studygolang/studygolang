// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"filter"
	"model"
	"net/http"
	"service"
	"strconv"
)

// 所有帖子（分页）
func TopicsHandler(rw http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}
	topics, _ := service.FindTopics(page, 0, "", "ctime DESC")
	// pageHtml := service.GetPageHtml(page, total)
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/topics.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"topics": topics})
}

// 所有节点（分页）
func NodesHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/nodes.html")
	filter.SetData(req, map[string]interface{}{"user": user, "nodes": model.AllNode})
}
