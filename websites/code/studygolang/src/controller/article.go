// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"service"
	"util"
)

// 网友文章列表页
// uri: /articles
func ArticlesHandler(rw http.ResponseWriter, req *http.Request) {
	lastId := req.FormValue("lastid")
	if lastId == "" {
		lastId = "0"
	}

	articles := service.FindArticles(lastId, "25")
	if articles == nil {
		// TODO:服务暂时不可用？
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != "0" {
		hasPrev = true
		prevId, _ = strconv.Atoi(lastId)
	}

	num := len(articles)

	if num > 20 {
		hasNext = true
		articles = articles[:20]
		nextId = articles[19].Id
	} else {
		nextId = articles[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"articles": articles, "activeArticles": "active", "page": pageInfo})
}

// 文章详细页
// uri: /articles/{id:[0-9]+}
func ArticleDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	article, err := service.FindArticleById(vars["id"])
	if err != nil {
		// TODO:
	}

	if article.Id == 0 {
		util.Redirect(rw, req, "/articles")
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/detail.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeArticles": "active", "article": article})
}
