// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package admin

import (
	"filter"
	"fmt"
	"html/template"
	"logger"
	"net/http"
	"service"
)

// 所有用户（分页）
func UserListHandler(rw http.ResponseWriter, req *http.Request) {
	curPage, limit := parsePage(req)

	users, total := service.FindUsersByPage(nil, curPage, limit)

	if users == nil {
		logger.Errorln("[UsersHandler]sql find error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"datalist":   users,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage + 1,
		"limit":      limit,
	}

	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/user/list.html,/template/admin/user/query.html")
	filter.SetData(req, data)
}

// 添加新用户表单页面
func NewUserHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	tpl, err := template.ParseFiles(ROOT+"/template/admin/common.html", ROOT+"/template/admin/newuser.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, user)
}

// 执行添加新用户（异步请求，返回json）
func AddUserHandler(rw http.ResponseWriter, req *http.Request) {
	// 入库
	errMsg, err := service.CreateUser(req.Form)
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error":"`, errMsg, `"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}

func ProfilerHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT+"/template/admin/common.html", ROOT+"/template/admin/profiler.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}
