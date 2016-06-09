// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type UserController struct{}

// 注册路由
func (self UserController) RegisterRoute(g *echo.Group) {
	g.GET("/user/:username", self.Home)
	g.GET("/users", self.ReadList)
	g.Match([]string{"GET", "POST"}, "/user/email/unsubscribe", self.EmailUnsub)
}

// Home 用户个人首页
func (UserController) Home(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/users")
	}

	topics := logic.DefaultTopic.FindRecent(5, user.Uid)

	resources := logic.DefaultResource.FindRecent(ctx, user.Uid)
	for _, resource := range resources {
		resource.CatName = logic.GetCategoryName(resource.Catid)
	}

	projects := logic.DefaultProject.FindRecent(ctx, user.Username)
	comments := logic.DefaultComment.FindRecent(ctx, user.Uid, -1, 5)

	return render(ctx, "user/profile.html", map[string]interface{}{"activeUsers": "active", "topics": topics, "resources": resources, "projects": projects, "comments": comments, "user": user})
}

// ReadList 会员列表
func (UserController) ReadList(ctx echo.Context) error {
	// 获取活跃会员
	activeUsers := logic.DefaultUser.FindActiveUsers(ctx, 36)
	// 获取最新加入会员
	newUsers := logic.DefaultUser.FindNewUsers(ctx, 36)
	// 获取会员总数
	total := logic.DefaultUser.Total()

	return render(ctx, "user/users.html", map[string]interface{}{"activeUsers": "active", "actives": activeUsers, "news": newUsers, "total": total})
}

// EmailUnsub 邮件订阅/退订页面
func (UserController) EmailUnsub(ctx echo.Context) error {
	token := ctx.FormValue("u")
	if token == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	// 校验 token 的合法性
	email := ctx.FormValue("email")
	user := logic.DefaultUser.FindOne(ctx, "email", email)
	if user.Email == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	realToken := logic.DefaultEmail.GenUnsubscribeToken(user)
	if token != realToken {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	if ctx.Request().Method() != "POST" {
		data := map[string]interface{}{
			"email":       email,
			"token":       token,
			"unsubscribe": user.Unsubscribe,
		}

		return render(ctx, "user/email_unsub.html", data)
	}

	logic.DefaultUser.EmailSubscribe(ctx, user.Uid, goutils.MustInt(ctx.FormValue("unsubscribe")))

	return success(ctx, nil)
}
