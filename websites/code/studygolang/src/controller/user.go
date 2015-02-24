// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"service"
	"util"
)

// 用户个人首页
// URI: /user/{username}
func UserHomeHandler(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	username := vars["username"]
	// 获取用户信息
	user := service.FindUserByUsername(username)

	if user == nil {
		util.Redirect(rw, req, "/users")
		return
	}

	topics := service.FindRecentTopics(user.Uid, "5")

	resources := service.FindUserRecentResources(user.Uid)
	resourceCats := make(map[int]string)
	for _, resource := range resources {
		resourceCats[resource.Catid] = service.GetCategoryName(resource.Catid)
	}

	projects := service.FindUserRecentProjects(user.Username)
	comments := service.FindRecentComments(user.Uid, -1, "5")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeUsers": "active", "topics": topics, "resources": resources, "resource_cats": resourceCats, "projects": projects, "comments": comments, "user": user})
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/user/profile.html")
}

// 会员列表
// URI: /users
func UsersHandler(rw http.ResponseWriter, req *http.Request) {
	// 获取活跃会员
	activeUsers := service.FindActiveUsers(0, 36)
	// 获取最新加入会员
	newUsers := service.FindNewUsers(0, 36)
	// 获取会员总数
	total := service.CountUsers()
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/user/users.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeUsers": "active", "actives": activeUsers, "news": newUsers, "total": total})
}

// 邮件订阅/退订页面
// URI: /user/email/unsubscribe{json:(|.json)}
func EmailUnsubHandler(rw http.ResponseWriter, req *http.Request) {
	token := req.FormValue("u")
	if token == "" {
		util.Redirect(rw, req, "/")
		return
	}

	// 校验 token 的合法性
	email := req.FormValue("email")
	user := service.FindUserByEmail(email)
	if user.Email == "" {
		util.Redirect(rw, req, "/")
		return
	}

	realToken := service.GenUnsubscribeToken(user.Username, user.Email)
	if token != realToken {
		util.Redirect(rw, req, "/")
		return
	}

	vars := mux.Vars(req)
	if req.Method != "POST" || vars["json"] == "" {
		filter.SetData(req, map[string]interface{}{
			"email":       email,
			"token":       token,
			"unsubscribe": user.Unsubscribe,
		})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/user/email_unsub.html")
		return
	}

	unsubscribe, _ := strconv.Atoi(req.PostFormValue("unsubscribe"))

	service.EmailSubscribe(user.Uid, unsubscribe)
	fmt.Fprint(rw, `{"ok": 1, "msg":"保存成功"}`)
}
