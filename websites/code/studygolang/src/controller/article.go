// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"github.com/studygolang/mux"
	"net/http"
	"service"
)

// 网友文章列表页
// uri: /articles
func ArticlesHandler(rw http.ResponseWriter, req *http.Request) {
	lastId := req.FormValue("lastid")
	if lastId == "" {
		lastId = "0"
	}

	articles := service.FindArticles(lastId, "21")
	if articles == nil {
		// TODO:服务暂时不可用？
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"articles": articles, "activeArticles": "active"})
}

// 文章详细页
// uri: /articles/{id:[0-9]+}
func ArticleDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	article, err := service.FindArticleById(vars["id"])
	if err != nil {
		// TODO:
	}
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/detail.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeArticles": "active", "article": article})
}
