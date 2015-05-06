// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"model"
	"service"
	"util"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	service.RegisterCommentObject(model.TYPE_CODE, service.CodeComment{})
	service.RegisterLikeObject(model.TYPE_CODE, service.CodeLike{})
}

// 代码片段列表页
// uri: /codes
func CodesHandler(rw http.ResponseWriter, req *http.Request) {

	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}

	codes, userMap, total := service.FindCodes(page)
	pageHtml := service.GetPageHtml(page, total, req.URL.Path)

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/codes/list.html")
	filter.SetData(req, map[string]interface{}{"activeCodes": "active", "codes": codes, "page": template.HTML(pageHtml), "userMap": userMap})
}

// 代码详细页
// uri: /codes/{id:[0-9]+}
func CodeDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	article, prevNext, err := service.FindArticlesById(vars["id"])
	if err != nil {
		util.Redirect(rw, req, "/articles")
		return
	}

	if article == nil || article.Id == 0 || article.Status == model.StatusOffline {
		util.Redirect(rw, req, "/articles")
		return
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

	// 为了阅读数即时看到
	article.Viewnum++

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/articles/detail.html,/template/common/comment.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeArticles": "active", "article": article, "prev": prevNext[0], "next": prevNext[1], "likeflag": likeFlag, "hadcollect": hadCollect})
}

// 分享代码
// uri: /codes/new{json:(|.json)}
func NewCodeHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title := req.PostFormValue("title")
	// 请求新建分享代码页面
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/codes/new.html")
		filter.SetData(req, map[string]interface{}{"activeCodes": "active"})
		return
	}

	errMsg := ""
	resForm := req.PostFormValue("form")
	if resForm == model.LinkForm {
		if req.PostFormValue("url") == "" {
			errMsg = "url不能为空"
		}
	} else {
		if req.PostFormValue("content") == "" {
			errMsg = "内容不能为空"
		}
	}
	if errMsg != "" {
		fmt.Fprint(rw, `{"ok": 0, "error":"`+errMsg+`"}`)
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishResource(user, req.PostForm)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误，请稍候再试！"}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}

// 修改資源
// uri: /codes/modify{json:(|.json)}
func ModifyCodeHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		util.Redirect(rw, req, "/resources")
		return
	}

	vars := mux.Vars(req)
	// 请求编辑資源页面
	if req.Method != "POST" || vars["json"] == "" {
		resource := service.FindResourceById(id)
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/new.html")
		filter.SetData(req, map[string]interface{}{"resource": resource, "activeResources": "active", "categories": service.AllCategory})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishResource(user, req.PostForm)
	if err != nil {
		if err == service.NotModifyAuthorityErr {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}
