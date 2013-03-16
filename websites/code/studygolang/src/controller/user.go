// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"filter"
	"github.com/studygolang/mux"
	"model"
	"net/http"
	"service"
)

// 用户个人首页
// URI: /user/{username}
func UserHomeHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/user/profile.html")
	vars := mux.Vars(req)
	username := vars["username"]
	// 获取用户信息
	user := service.FindUserByUsername(username)
	if user != nil {
		topics := service.FindRecentTopics(user.Uid)
		comments := service.FindRecentComments(user.Uid, model.TYPE_TOPIC)
		replies := service.FindRecentReplies(comments)
		// 设置模板数据
		filter.SetData(req, map[string]interface{}{"activeUsers": "active", "topics": topics, "replies": replies, "user": user})
	}
}

// 会员列表
// URI: /users
func UsersHandler(rw http.ResponseWriter, req *http.Request) {
	// 获取活跃会员
	activeUsers := service.FindActiveUsers(0, 30)
	// 获取最新加入会员
	newUsers := service.FindNewUsers(0, 30)
	// 获取会员总数
	total := service.CountUsers()
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/user/users.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeUsers": "active", "actives": activeUsers, "news": newUsers, "total": total})
}
