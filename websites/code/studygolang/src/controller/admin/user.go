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
	"net/http"
	"service"
)

// 所有用户（分页）
func UsersHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	users, err := service.FindUsers()
	if err != nil {
		// TODO:
	}
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/users.html")
	filter.SetData(req, map[string]interface{}{"user": user, "users": users})
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
