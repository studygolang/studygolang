// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"filter"
	"logger"
	"model"
	"net/http"
	"service"
	"strconv"
)

// 所有文章（分页）
// /admin/crawl/article/list
func ArticleListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	articles, total := service.FindArticleByPage(nil, curPage, limit)

	if articles == nil {
		logger.Errorln("[ArticleListHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   articles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/article/list.html,/template/admin/article/query.html")
	filter.SetData(req, data)
}

// /admin/crawl/article/query.html
func ArticleQueryHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	conds := parseConds(req, []string{"domain"})

	articles, total := service.FindArticleByPage(conds, curPage, limit)

	if articles == nil {
		logger.Errorln("[ArticleQueryHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   articles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/article/query.html")
	filter.SetData(req, data)
}

func ModifyArticleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	if req.PostFormValue("submit") == "1" {
		user, _ := filter.CurrentUser(req)

		errMsg, err := service.ModifyArticle(user, req.PostForm)
		if err != nil {
			data["ok"] = 0
			data["error"] = errMsg
		} else {
			data["ok"] = 1
			data["msg"] = "修改成功"
		}
	} else {
		article, err := service.FindArticleById(req.FormValue("id"))

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 设置内容模板
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/article/modify.html")
		data["article"] = article
		data["statusSlice"] = model.StatusSlice
		data["langSlice"] = model.LangSlice
	}

	filter.SetData(req, data)
}

func DelArticleHandler(rw http.ResponseWriter, req *http.Request) {
	var data = make(map[string]interface{})

	id := req.FormValue("id")

	if _, err := strconv.Atoi(id); err != nil {
		data["ok"] = 0
		data["error"] = "id不是整型"

		filter.SetData(req, data)
		return
	}

	if err := service.DelArticle(id); err != nil {
		data["ok"] = 0
		data["error"] = "删除失败！"
	} else {
		data["ok"] = 1
		data["msg"] = "删除成功！"
	}

	filter.SetData(req, data)
}
