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
	"service"
)

// search
// uri: /search
func SearchHandler(rw http.ResponseWriter, req *http.Request) {
	q := req.FormValue("q")
	field := req.FormValue("f")
	p, err := strconv.Atoi(req.FormValue("p"))
	if err != nil {
		p = 1
	}

	rows := 20

	pageHtml := ""
	respBody, err := service.DoSearch(q, field, (p-1)*rows, rows)
	if err == nil {
		pageHtml = service.GenPageHtml(p, rows, respBody.NumFound, "/search?q="+q+"&f="+field)
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/search.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"respBody": respBody, "q": q, "f": field, "pageHtml": pageHtml})
}
