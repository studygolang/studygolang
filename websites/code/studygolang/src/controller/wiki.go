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
	"net/http"
	"service"
)

// 创建wiki页
// uri: /wiki/new{json:(|.json)}
func NewWikiPageHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title := req.FormValue("title")
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/wiki/new.html")
		filter.SetData(req, map[string]interface{}{"activeWiki": "active"})
		return
	}
	user, _ := filter.CurrentUser(req)
	// 入库
	ok := service.CreateWiki(user["uid"].(int), req.Form)
	if !ok {
		fmt.Fprint(rw, `{"errno": 1, "error":"服务器内部错误，请稍候再试！"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "data":{"uri":"`+req.FormValue("uri")+`"}}`)
}

// 展示wiki页
// uri: /wiki/{uri}
func WikiContentHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uri := vars["uri"]
	wiki := service.FindWiki(uri)
	if wiki == nil {
		NotFoundHandler(rw, req)
		return
	}
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/wiki/content.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeWiki": "active", "wiki": wiki})
}

// 获得wiki列表
// uri: /wiki
func WikisHandler(rw http.ResponseWriter, req *http.Request) {
	wikiList := service.FindWikiList()
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/wiki/list.html")
	filter.SetData(req, map[string]interface{}{"activeWiki": "active", "wikis": wikiList})
}
