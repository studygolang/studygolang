// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
)

// subrouter.HandleFunc("/crawl/article/list", admin.ArticleListHandler)
// 	subrouter.HandleFunc("/crawl/article/query.html", admin.ArticleQueryHandler)
// 	subrouter.HandleFunc("/crawl/article/modify", admin.ModifyArticleHandler)
// 	subrouter.HandleFunc("/crawl/article/new", admin.CrawlArticleHandler)
// 	subrouter.HandleFunc("/crawl/article/del", admin.DelArticleHandler)

type ArticleController struct{}

// 注册路由
func (self ArticleController) RegisterRoute(g *echo.Group) {
	g.Get("/crawl/article/list", echo.HandlerFunc(self.ArticleList))
	g.Post("/crawl/article/query.html", echo.HandlerFunc(self.ArticleQuery))
}

// ArticleList 所有文章（分页）
func (ArticleController) ArticleList(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	articles, total := logic.DefaultArticle.FindArticleByPage(ctx, nil, curPage, limit)

	if articles == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   articles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "article/list.html,article/query.html", data)
}

// ArticleQuery
func (ArticleController) ArticleQuery(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	conds := parseConds(ctx, []string{"id", "domain", "title"})

	articles, total := logic.DefaultArticle.FindArticleByPage(ctx, conds, curPage, limit)

	if articles == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   articles,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return renderQuery(ctx, "article/query.html", data)
}

// // /admin/crawl/article/new
// func CrawlArticleHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	if req.PostFormValue("submit") == "1" {
// 		urls := strings.Split(req.PostFormValue("urls"), "\n")

// 		var errMsg string
// 		for _, articleUrl := range urls {
// 			_, err := service.ParseArticle(strings.TrimSpace(articleUrl), false)

// 			if err != nil {
// 				errMsg = err.Error()
// 			}
// 		}

// 		if errMsg != "" {
// 			data["ok"] = 0
// 			data["error"] = errMsg
// 		} else {
// 			data["ok"] = 1
// 			data["msg"] = "添加成功"
// 		}
// 	} else {

// 		// 设置内容模板
// 		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/article/new.html")
// 	}

// 	filter.SetData(req, data)
// }

// func ModifyArticleHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	if req.PostFormValue("submit") == "1" {
// 		user, _ := filter.CurrentUser(req)

// 		errMsg, err := service.ModifyArticle(user, req.PostForm)
// 		if err != nil {
// 			data["ok"] = 0
// 			data["error"] = errMsg
// 		} else {
// 			data["ok"] = 1
// 			data["msg"] = "修改成功"
// 		}
// 	} else {
// 		article, err := service.FindArticleById(req.FormValue("id"))

// 		if err != nil {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// 设置内容模板
// 		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/article/modify.html")
// 		data["article"] = article
// 		data["statusSlice"] = model.StatusSlice
// 		data["langSlice"] = model.LangSlice
// 	}

// 	filter.SetData(req, data)
// }

// func DelArticleHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	id := req.FormValue("id")

// 	if _, err := strconv.Atoi(id); err != nil {
// 		data["ok"] = 0
// 		data["error"] = "id不是整型"

// 		filter.SetData(req, data)
// 		return
// 	}

// 	if err := service.DelArticle(id); err != nil {
// 		data["ok"] = 0
// 		data["error"] = "删除失败！"
// 	} else {
// 		data["ok"] = 1
// 		data["msg"] = "删除成功！"
// 	}

// 	filter.SetData(req, data)
// }
