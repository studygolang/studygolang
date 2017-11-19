// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	. "http"
	"http/middleware"
	"logic"
	"model"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/echoutils"
)

type GCTTController struct{}

// 注册路由
func (self GCTTController) RegisterRoute(g *echo.Group) {
	g.Get("/gctt", self.Index)
	g.Get("/gctt-list", self.UserList)
	g.Get("/gctt/:username", self.User)
	g.Get("/gctt-apply", self.Apply, middleware.NeedLogin())
	g.Match([]string{"GET", "POST"}, "/gctt-new", self.Create, middleware.NeedLogin())
}

func (self GCTTController) Index(ctx echo.Context) error {
	return Render(ctx, "gctt/index.html", map[string]interface{}{})
}

// Apply 申请成为译者
func (GCTTController) Apply(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	gcttUser := logic.DefaultGCTT.FindTranslator(ctx, me)
	if gcttUser.Id > 0 {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	// 是否绑定了 github 账号
	var githubUser *model.BindUser
	bindUsers := logic.DefaultUser.FindBindUsers(ctx, me.Uid)
	for _, bindUser := range bindUsers {
		if bindUser.Type == model.BindTypeGithub {
			githubUser = bindUser
			break
		}
	}

	// 如果已经绑定，查看是否之前已经是译者
	if githubUser != nil {
		gcttUser = logic.DefaultGCTT.FindOne(ctx, githubUser.Username)
		logic.DefaultGCTT.BindUser(ctx, gcttUser, me.Uid, githubUser)
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	return render(ctx, "gctt/apply.html", map[string]interface{}{
		"activeGCTT":  "active",
		"github_user": githubUser,
	})
}

// Create 发布新译文
func (GCTTController) Create(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	gcttUser := logic.DefaultGCTT.FindTranslator(ctx, me)

	title := ctx.FormValue("title")
	if title == "" || ctx.Request().Method() != "POST" {
		return render(ctx, "gctt/new.html", map[string]interface{}{
			"activeGCTT": "active",
			"gctt_user":  gcttUser,
		})
	}

	if ctx.FormValue("content") == "" {
		return fail(ctx, 1, "内容不能为空")
	}

	if gcttUser == nil {
		return fail(ctx, 2, "不允许发布！")
	}

	id, err := logic.DefaultArticle.Publish(echoutils.WrapEchoContext(ctx), me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 3, "内部服务错误")
	}

	return success(ctx, map[string]interface{}{"id": id})
}

func (GCTTController) User(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	gcttUser := logic.DefaultGCTT.FindOne(ctx, username)
	if gcttUser.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	return render(ctx, "gctt/user-info.html", map[string]interface{}{
		"gctt_user": gcttUser,
	})
}

func (GCTTController) UserList(ctx echo.Context) error {
	return render(ctx, "gctt/user-list.html", map[string]interface{}{})
}
