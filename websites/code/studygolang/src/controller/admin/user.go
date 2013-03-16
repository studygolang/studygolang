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
	"util"
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
	rules := map[string]map[string]map[string]string{
		"username": {
			"require": {"error": "用户名不能为空！"},
			"length":  {"range": "4,20", "error": "用户名长度必须在%d个字符和%d个字符之间"},
		},
		"email": {
			"require": {"error": "邮箱不能为空！"},
			"email":   {"error": "邮箱格式不正确！"},
		},
		"passwd": {
			"require": {"error": "密码不能为空！"},
			"length":  {"range": "6,32", "error": "密码长度必须在%d个字符和%d个字符之间"},
		},
		"pass2": {
			"require": {"error": "确认密码不能为空！"},
			"compare": {"field": "passwd", "rule": "=", "error": "两次密码不一致"},
		},
	}
	req.ParseForm()
	errMsg := util.Validate(req.Form, rules)

	header := rw.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	if errMsg != "" {
		fmt.Fprint(rw, `{"errno": 1,"error":"`, errMsg, `"}`)
		return
	}

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
