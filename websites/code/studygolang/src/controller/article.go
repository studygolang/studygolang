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
	"model"
	"service"
	"util"
)

const limit = 20

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	service.RegisterCommentObject(model.TYPE_ARTICLE, service.ArticleComment{})
	service.RegisterLikeObject(model.TYPE_ARTICLE, service.ArticleLike{})
}

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

	num := len(articles)

	if num == 0 {
		if lastId == "0" {
			util.Redirect(rw, req, "/")
		} else {
			util.Redirect(rw, req, "/articles")
		}
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != "0" {
		prevId, _ = strconv.Atoi(lastId)

		// 避免因为文章下线，导致判断错误（所以 > 5）
		if prevId-articles[0].Id > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		articles = articles[:limit]
		nextId = articles[limit-1].Id
	} else {
		nextId = articles[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	// 获取当前用户喜欢对象信息
	user, ok := filter.CurrentUser(req)
	var likeFlags map[int]int
	if ok {
		uid := user["uid"].(int)
		likeFlags, _ = service.FindUserLikeObjects(uid, model.TYPE_ARTICLE, articles[0].Id, nextId)
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"articles": articles, "activeArticles": "active", "page": pageInfo, "likeflags": likeFlags})
}

// 文章详细页
// uri: /articles/{id:[0-9]+}
func ArticleDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	article, prevNext, err := service.FindArticlesById(vars["id"])
	if err != nil {
		util.Redirect(rw, req, "/articles")
	}

	if article == nil || article.Id == 0 || article.Status == model.StatusOffline {
		util.Redirect(rw, req, "/articles")
	}

	likeFlag := 0
	hadCollect := 0
	user, ok := filter.CurrentUser(req)
	if ok {
		uid := user["uid"].(int)
		likeFlag = service.HadLike(uid, article.Id, model.TYPE_ARTICLE)
		hadCollect = service.HadFavorite(uid, article.Id, model.TYPE_ARTICLE)
	}

	service.Views.Incr(req, model.TYPE_ARTICLE, article.Id)

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/detail.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeArticles": "active", "article": article, "prev": prevNext[0], "next": prevNext[1], "likeflag": likeFlag, "hadcollect": hadCollect})
}
