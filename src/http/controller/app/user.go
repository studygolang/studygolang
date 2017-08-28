// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"logic"
	"model"

	"github.com/labstack/echo"

	. "http"
	. "http/internal/helper"
)

type UserController struct{}

// 注册路由
func (self UserController) RegisterRoute(g *echo.Group) {
	g.GET("/user/center", self.Center)
	g.GET("/user/me", self.Me)
	g.POST("/user/modify", self.Modify)
	g.POST("/user/login", self.Login)
}

// Center 用户自己个人中心
func (UserController) Center(ctx echo.Context) error {
	if user, ok := ctx.Get("user").(*model.Me); ok {
		data := map[string]interface{}{
			"user": user,
		}
		return success(ctx, data)
	}

	return success(ctx, nil)
}

// Me 用户信息
func (UserController) Me(ctx echo.Context) error {
	if me, ok := ctx.Get("user").(*model.Me); ok {
		user := logic.DefaultUser.FindOne(ctx, "uid", me.Uid)
		return success(ctx, map[string]interface{}{
			"user":            user,
			"default_avatars": logic.DefaultAvatars,
		})
	}

	return success(ctx, nil)
}

func (UserController) Login(ctx echo.Context) error {
	if _, ok := ctx.Get("user").(*model.Me); ok {
		return success(ctx, nil)
	}

	username := ctx.FormValue("username")
	if username == "" {
		return fail(ctx, "用户名为空")
	}

	// 处理用户登录
	passwd := ctx.FormValue("passwd")
	userLogin, err := logic.DefaultUser.Login(ctx, username, passwd)
	if err != nil {
		return fail(ctx, err.Error())
	}

	data := map[string]interface{}{
		"token":    GenToken(userLogin.Uid),
		"uid":      userLogin.Uid,
		"username": userLogin.Username,
	}
	return success(ctx, data)
}

func (UserController) Modify(ctx echo.Context) error {
	me, ok := ctx.Get("user").(*model.Me)
	if !ok {
		return fail(ctx, "请先登录", NeedReLoginCode)
	}

	// 更新信息
	errMsg, err := logic.DefaultUser.Update(ctx, me, ctx.Request().FormParams())
	if err != nil {
		return fail(ctx, errMsg)
	}

	email := ctx.FormValue("email")
	if me.Email != email {
		isHttps := CheckIsHttps(ctx)
		go logic.DefaultEmail.SendActivateMail(email, RegActivateCode.GenUUID(email), isHttps)
	}

	return success(ctx, nil)
}
