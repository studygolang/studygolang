// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"model"
	"net/http"
	"service"
	"util"
)

// 在需要评论且要回调的地方注册评论对象
func init() {
	// 注册评论对象
	service.RegisterCommentObject("resource", service.ResourceComment{})
}

// 资源索引页
// uri: /resources
func ResIndexHandler(rw http.ResponseWriter, req *http.Request) {
	util.Redirect(rw, req, "/resources/cat/1")
}

// 某个分类的资源列表
// uri: /resources/cat/{catid:[0-9]+}
func CatResourcesHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	catid := vars["catid"]
	resources := service.FindResourcesByCatid(catid)
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/index.html")
	filter.SetData(req, map[string]interface{}{"activeResources": "active", "resources": resources, "categories": service.AllCategory, "curCatid": catid})
}

// 某个资源详细页
// uri: /resources/{id:[0-9]+}
func ResourceDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	resource, comments := service.FindResource(vars["id"])
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/detail.html")
	filter.SetData(req, map[string]interface{}{"activeResources": "active", "resource": resource, "comments": comments})
}

// 发布新资源
// uri: /resources/new{json:(|.json)}
func NewResourceHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title := req.FormValue("title")
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/new.html")
		filter.SetData(req, map[string]interface{}{"activeResources": "active", "categories": service.AllCategory})
		return
	}
	errMsg := ""
	resForm := req.FormValue("form")
	if resForm == model.LinkForm {
		if req.FormValue("url") == "" {
			errMsg = "url不能为空"
		}
	} else {
		if req.FormValue("content") == "" {
			errMsg = "内容不能为空"
		}
	}
	if errMsg != "" {
		fmt.Fprint(rw, `{"errno": 1, "error":"`+errMsg+`"}`)
		return
	}
	user, _ := filter.CurrentUser(req)
	// 入库
	ok := service.PublishResource(user["uid"].(int), req.Form)
	if !ok {
		fmt.Fprint(rw, `{"errno": 1, "error":"服务器内部错误，请稍候再试！"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "data":""}`)
}
